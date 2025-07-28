package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/Snoop-Duck/ToDoList/internal"
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

type Repository interface {
	SaveUser(user users.User) error
	GetUser(login string) (users.User, error)
	DeleteUser(userID string) error
	GetAllUsers() ([]users.User, error)
	GetUserID(userID string) (users.User, error)
	UpdateUserID(userID string, user users.User) error
	Close() error
}

type RepositoryNote interface {
	AddNote(note notes.Note) error
	GetNotes() ([]notes.Note, error)
	GetNoteID(noteID string) (notes.Note, error)
	DeleteNote(noteID string) error
	UpdateNote(noteID string, note notes.Note) error
}

type NotesAPI struct {
	cfg       *internal.Config
	httpServe *http.Server
	repo      Repository
	repoNote  RepositoryNote
	log       zerolog.Logger
	testMode  bool
}

func New(cfg *internal.Config, repo Repository, repoNote RepositoryNote) *NotesAPI {
	var log zerolog.Logger
	if cfg != nil {
		log = logger.Get(cfg.Debug)
	} else {
		log = zerolog.Nop() // для тестов
	}
	log.Debug().Msg("configure Notes API server")
	httpServe := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	notesAPI := NotesAPI{
		httpServe: &httpServe,
		cfg:       cfg,
		repo:      repo,
		repoNote:  repoNote,
		log:       log,
	}
	notesAPI.configRoutes()
	return &notesAPI
}

func (nApi *NotesAPI) Run() error {
	nApi.log.Info().Msgf("notes API started on %s", nApi.httpServe.Addr)
	return nApi.httpServe.ListenAndServe()
}

func (nApi *NotesAPI) Stop(ctx context.Context) error {
	return nApi.httpServe.Shutdown(ctx)
}

func (nApi *NotesAPI) configRoutes() {
	nApi.log.Debug().Msg("configure routes")
	router := gin.Default()

	router.Use(gzip.Gzip(
		gzip.BestSpeed,
		gzip.WithDecompressFn(gzip.DefaultDecompressHandle),
	))

	router.Use(gzip.Gzip(
		gzip.BestSpeed,
		gzip.WithExcludedExtensions([]string{".png", ".gif", ".jpeg", ".jpg"}),
		gzip.WithExcludedPaths([]string{"/metrics"}),
	))

	router.Use(func(c *gin.Context) {
		contentType := c.Writer.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") &&
			!strings.Contains(contentType, "text/html") {
			c.Header("Content-Encoding", "identity")
		}
		c.Next()
	})

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
		notes.GET("/list/:id", nApi.getNoteID)
		notes.POST("/add", nApi.createNote)
		notes.PUT("/upd/:id", nApi.updateNote)
		notes.DELETE("/del/:id", nApi.deleteNote)
	}
	nApi.httpServe.Handler = router
}

func (nApi *NotesAPI) JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if nApi.testMode {
			ctx.Set("uid", "test-user")
			ctx.Next()
			return
		}

		token := ctx.GetHeader("Authorization")
		if token == `` {
			nApi.log.Error().Msg("token not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		uid, err := validateJwtToken(token)
		if err != nil {
			nApi.log.Error().Err(err).Msg("failed to validate token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
			return
		}
		nApi.log.Debug().Str("uid", uid).Msg("user was authorized")
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
