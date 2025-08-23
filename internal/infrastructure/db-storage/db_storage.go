package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog"

	// Импорт драйвера PostgreSQL для миграций.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Импорт драйвера PostgreSQL для database/sql.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	deleteBatchSize = 10
	deleteChanSize  = 100
)

type DBStorage struct {
	db         *pgxpool.Pool
	deleteChan chan struct{}
	stopChan   chan struct{}
	log        zerolog.Logger
}

func New(ctx context.Context, addr string) (*DBStorage, error) {
	pool, err := pgxpool.New(ctx, addr)
	if err != nil {
		return nil, err
	}

	storage := &DBStorage{
		db:         pool,
		deleteChan: make(chan struct{}, deleteChanSize),
		stopChan:   make(chan struct{}),
		log:        logger.Get(),
	}

	go storage.startBatchDeleter()
	return storage, nil
}

func (db *DBStorage) Close() error {
	close(db.stopChan)
	db.db.Close()
	return nil
}

func ApplyMigrations(addr string) error {
	migrationPath := "file://migrations"
	m, err := migrate.New(migrationPath, addr)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func (db *DBStorage) startBatchDeleter() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-db.stopChan:
			if err := db.processBatchDeletion(); err != nil {
				db.log.Error().Err(err).Msg("failed to process batch deletion on shutdown: %v")
			}
			return
		case <-ticker.C:
			if err := db.processBatchDeletion(); err != nil {
				db.log.Error().Err(err).Msg("failed to process batch deletion: %v")
			}
		case <-db.deleteChan:
			if len(db.deleteChan) >= deleteBatchSize-1 {
				if err := db.processBatchDeletion(); err != nil {
					db.log.Error().Err(err).Msg("failed to process batch deletion: %v")
				}
			}
		}
	}
}

func (db *DBStorage) processBatchDeletion() error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			db.log.Error().Err(rollbackErr).Msg("failed to rollback transaction: %v")
		}
	}()

	_, err = tx.Exec(ctx, `
		DELETE FROM notes 
		WHERE nid IN (
			SELECT nid FROM notes 
			WHERE deleted = true 
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)`, deleteBatchSize)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
