package main

import (
	"fmt"
	"main/ToDoList/internal"
	inmemorynotes "main/ToDoList/internal/infrastructure/notes"
	inmemory "main/ToDoList/internal/infrastructure/users"
	"main/ToDoList/internal/server"
)

func main() {
	cfg := internal.ReadConfig()
	fmt.Printf("Host: %s\nPort: %d\n", cfg.Host, cfg.Port)
	noteRepo := inmemorynotes.New()
	userRepo := inmemory.New()
	notesAPI := server.New(cfg, userRepo, noteRepo)
	notesAPI.Run()
}
