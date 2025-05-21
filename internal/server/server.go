package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/Snoop-Duck/ToDoList/internal"
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

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
	log := logger.Get()
	log.Debug().Msg("configure Notes API server")
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
	log := logger.Get()
	log.Info().Msgf("notes API started on %s", nApi.httpServe.Addr)
	return nApi.httpServe.ListenAndServe()
}

func (nApi *NotesAPI) Stop(ctx context.Context) error {
	return nApi.httpServe.Shutdown(ctx)
}

func (nApi *NotesAPI) configRoutes() {
	log := logger.Get()
	log.Debug().Msg("configure routes")
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
		notes.GET("/list", nApi.JWTMiddleware(), nApi.getNotes)
		notes.GET("list/:id", nApi.getNoteID)
		notes.POST("/add", nApi.createNote)
		notes.PUT("upd/:id", nApi.updateNote)
		notes.DELETE("del/:id", nApi.deleteNote)
	}
	nApi.httpServe.Handler = router
}

func (nApi *NotesAPI) JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.Get()
		token := ctx.GetHeader("Authorization")
		if token == `` {
			log.Error().Msg("token not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uid, err := validateJwtToken(token)
		if err != nil {
			log.Error().Err(err).Msg("failed to validate token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
			return
		}
		log.Debug().Str("uid", uid).Msg("user was authorized")
		ctx.Set("uid", uid)
		ctx.Next()
	}
}

func jwtToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   uid,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return ``, err
	}
	return tokenString, nil
}

func validateJwtToken(tokenString string) (string, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return ``, err
	}

	if !token.Valid {
		return ``, fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}
