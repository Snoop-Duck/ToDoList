package dbstorage

import (
	"context"

	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	db *pgx.Conn
}

func New(ctx context.Context, addr string) (*DBStorage, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &DBStorage{db: conn}, nil
}

func (db *DBStorage) Close() error {
	return db.db.Close(context.Background())
}

func AppyMigrations(addr string) error {
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
