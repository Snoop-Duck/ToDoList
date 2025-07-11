package dbstorage

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
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

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (db *DBStorage) startBatchDeleter() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-db.stopChan:
			db.processBatchDeletion()
			return
		case <-ticker.C:
			db.processBatchDeletion()
		case <-db.deleteChan:
			if len(db.deleteChan) >= deleteBatchSize-1 {
				db.processBatchDeletion()
			}
		}
	}
}

func (db *DBStorage) processBatchDeletion() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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
