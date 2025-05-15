package server

import (
	"context"
	"fmt"
	"github.com/Snoop-Duck/ToDoList/internal"
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	SaveUser(user users.User) error
	GetUser(login string) (users.User, error)
	DeleteUser(userID string) error
	GetAllUsers() (map[string]users.User, error)
	GetUserID(userID string) (users.User, error)
	UpdateUserID(userID string, user users.User) error
}

type RepositoryNote interface {
	AddNote(note notes.Note) error
	GetNotes() (map[string]notes.Note, error)
	GetNoteID(noteID string) (notes.Note, error)
	DeleteNote(noteID string) error
	UpdateNote(noteID string, note notes.Note) error
}

type NotesAPI struct {
	cfg       *internal.Config
	httpServe *http.Server
	repo      Repository
	repoNote  RepositoryNote
}

func New(cfg *internal.Config, repo Repository, repoNote RepositoryNote) *NotesAPI {
	httpServe := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	notesAPI := NotesAPI{
		httpServe: &httpServe,
		cfg:       cfg,
		repo:      repo,
		repoNote:  repoNote,
	}
	notesAPI.configRoutes()
	return &notesAPI
}

func (nApi *NotesAPI) Run() error {
	return nApi.httpServe.ListenAndServe()
}

func (nApi *NotesAPI) Stop(ctx context.Context) error {
	return nApi.httpServe.Shutdown(ctx)
}

func (nApi *NotesAPI) configRoutes() {
	router := gin.Default()
	router.GET("/")
	users := router.Group("/users")
	{
		users.GET("/profile", nApi.getUsers)
		users.GET("/profile/:id", nApi.getUserID)
		users.POST("/register", nApi.register)
		users.POST("/login", nApi.login)
		users.PUT("/upd/:id", nApi.updateUserID)
		users.DELETE("/del/:id", nApi.deleteUser)
	}
	notes := router.Group("/notes")
	{
		notes.GET("/list", nApi.getNotes)
		notes.GET("list/:id", nApi.getNoteID)
		notes.POST("/add", nApi.createNote)
		notes.PUT("upd/:id", nApi.updateNote)
		notes.DELETE("del/:id", nApi.deleteNote)
	}
	nApi.httpServe.Handler = router
}
