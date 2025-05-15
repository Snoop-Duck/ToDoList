package main

import (
	"fmt"

	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/users"

	"github.com/Snoop-Duck/ToDoList/internal/server"

	inmemorynotes "github.com/Snoop-Duck/ToDoList/internal/infrastructure/notes"

	"github.com/Snoop-Duck/ToDoList/internal"
)

func main() {
	cfg := internal.ReadConfig()
	fmt.Printf("Host: %s\nPort: %d\n", cfg.Host, cfg.Port)
	noteRepo := inmemorynotes.New()
	userRepo := inmemory.New()
	notesAPI := server.New(cfg, userRepo, noteRepo)
	notesAPI.Run()
}
