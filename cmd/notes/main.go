package main

import (
	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/users"

	"github.com/Snoop-Duck/ToDoList/internal/server"

	inmemorynotes "github.com/Snoop-Duck/ToDoList/internal/infrastructure/notes"

	"github.com/Snoop-Duck/ToDoList/internal"
	logger "github.com/Snoop-Duck/ToDoList/pkg"
)

func main() {
	cfg := internal.ReadConfig()

	log := logger.Get(cfg.Debug)

	log.Info().Msg("service starting")

	noteRepo := inmemorynotes.New()
	userRepo := inmemory.New()
	notesAPI := server.New(cfg, userRepo, noteRepo)
	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("fatal running server")
	}
}
