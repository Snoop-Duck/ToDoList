package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type NotesAPITest struct {
	*NotesAPI
}

func NewTestNotesAPI(repoNote RepositoryNote) *NotesAPITest {
	return &NotesAPITest{
		NotesAPI: &NotesAPI{
			log:      zerolog.Nop(), // Используем Nop логгер в тестах
			repoNote: repoNote,
			testMode: true, // Ключевое изменение - включаем тестовый режим
		},
	}
}

func TestGetNotes(t *testing.T) {
	testNotes := []notes.Note{
		{NID: "1", Title: "Test Note 1", Status: notes.New},
		{NID: "2", Title: "Test Note 2", Status: notes.Active},
	}

	mockRepo := new(mocks.RepositoryNote)
	mockRepo.On("GetNotes").Return(testNotes, nil)

	api := NewTestNotesAPI(mockRepo)

	r := gin.New()
	r.GET("/notes/list", api.JWTMiddleware(), api.getNotes) // Используем встроенный middleware

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.R().Get(ts.URL + "/notes/list")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode())
	assert.Contains(t, resp.String(), "Test Note 1")
	assert.Contains(t, resp.String(), "Test Note 2")
	mockRepo.AssertExpectations(t)
}

func TestCreateNote(t *testing.T) {
	mockRepo := new(mocks.RepositoryNote)
	mockRepo.On("AddNote", mock.AnythingOfType("notes.Note")).Return(nil)

	api := NewTestNotesAPI(mockRepo)

	r := gin.New()
	r.POST("/notes", api.createNote)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"title":"New Note","status":0,"uid":"user1"}`).
		Post(ts.URL + "/notes")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	mockRepo.AssertExpectations(t)
}

func TestGetNoteByID(t *testing.T) {
	testNote := notes.Note{
		NID:    "123",
		Title:  "Test Note",
		Status: notes.Active,
	}

	mockRepo := new(mocks.RepositoryNote)
	mockRepo.On("GetNoteID", "123").Return(testNote, nil)

	api := NewTestNotesAPI(mockRepo)

	r := gin.New()
	r.GET("/notes/:id", api.getNoteID)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.R().Get(ts.URL + "/notes/123")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode())
	assert.Contains(t, resp.String(), "Test Note")
	mockRepo.AssertExpectations(t)
}

func TestUpdateNote(t *testing.T) {
	mockRepo := new(mocks.RepositoryNote)
	mockRepo.On("UpdateNote", "123", mock.AnythingOfType("notes.Note")).Return(nil)

	api := NewTestNotesAPI(mockRepo)

	r := gin.New()
	r.PUT("/notes/:id", api.updateNote)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.R().
		SetBody(`{"title":"Updated Note","status":1}`).
		Put(ts.URL + "/notes/123")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	mockRepo.AssertExpectations(t)
}

func TestDeleteNote(t *testing.T) {
	mockRepo := new(mocks.RepositoryNote)
	mockRepo.On("DeleteNote", "123").Return(nil)

	api := NewTestNotesAPI(mockRepo)

	r := gin.New()
	r.DELETE("/notes/:id", api.deleteNote)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.R().Delete(ts.URL + "/notes/123")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	mockRepo.AssertExpectations(t)
}
