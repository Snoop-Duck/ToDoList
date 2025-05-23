package inmemory

import (
	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
)

type InMemoryNotes struct {
	noteStorage map[string]notes.Note
}

var emtyUser = users.User{}

type InMemoryUsers struct {
	userStorage map[string]users.User
}

func NewNotes() *InMemoryNotes {
	return &InMemoryNotes{
		noteStorage: make(map[string]notes.Note),
	}
}

func NewUsers() *InMemoryUsers {
	log := logger.Get()
	log.Debug().Msg("create in memory storage")
	return &InMemoryUsers{
		userStorage: make(map[string]users.User),
	}
}
