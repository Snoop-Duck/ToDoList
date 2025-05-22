package main

import (
	"context"

	"github.com/Snoop-Duck/ToDoList/internal/server"

	dbstorage "github.com/Snoop-Duck/ToDoList/internal/infrastructure/db-storage"
	inmemorynotes "github.com/Snoop-Duck/ToDoList/internal/infrastructure/notes"
	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/users"

	"github.com/Snoop-Duck/ToDoList/internal"
	logger "github.com/Snoop-Duck/ToDoList/pkg"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := internal.ReadConfig()

	log := logger.Get(cfg.Debug)

	log.Info().Msg("service starting")

	var repoUser server.Repository
	var err error
	noteRepo := inmemorynotes.New()
	repoUser, err = dbstorage.New(context.Background(), "postgres://user:password@localhost:5432/notes?sslmode=disable")
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to db. Use in memory storage")
		repoUser = inmemory.New()
	}
	if err = dbstorage.AppyMigrations("postgres://user:password@localhost:5432/notes?sslmode=disable"); err != nil {
		log.Warn().Err(err).Msg("failed to apply migrations. Use in memory storage")
		repoUser.Close()
		repoUser = inmemory.New()
	}

	notesAPI := server.New(cfg, repoUser, noteRepo)
	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("fatal running server")
	}
}
