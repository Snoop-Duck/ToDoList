package main

import (
	"context"

	"github.com/Snoop-Duck/ToDoList/internal/server"

	dbstorage "github.com/Snoop-Duck/ToDoList/internal/infrastructure/db-storage"
	inmemorynotes "github.com/Snoop-Duck/ToDoList/internal/infrastructure/notes"
	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/users"

	"github.com/Snoop-Duck/ToDoList/internal"
	logger "github.com/Snoop-Duck/ToDoList/pkg"
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
	notesAPI := server.New(cfg, repoUser, noteRepo)
	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("fatal running server")
	}
}
