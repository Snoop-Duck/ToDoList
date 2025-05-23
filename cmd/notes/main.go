package main

import (
	"context"

	"github.com/Snoop-Duck/ToDoList/internal/server"

	dbstorage "github.com/Snoop-Duck/ToDoList/internal/infrastructure/db-storage"
	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/in-memory"

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
	repoNote := inmemory.NewNotes()
	repoUser, err = dbstorage.New(context.Background(), "postgres://user:password@localhost:5432/notes?sslmode=disable")
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to db. Use in memory storage")
		repoUser = inmemory.NewUsers()
	}
	if err = dbstorage.AppyMigrations("postgres://user:password@localhost:5432/notes?sslmode=disable"); err != nil {
		log.Warn().Err(err).Msg("failed to apply migrations. Use in memory storage")
		repoUser.Close()
		repoUser = inmemory.NewUsers()
	}

	notesAPI := server.New(cfg, repoUser, repoNote)
	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("fatal running server")
	}
}
