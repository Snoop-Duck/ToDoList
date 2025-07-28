package note

import (
	"errors"
	"testing"

	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/services/note/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNoteService_CreateNote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)
		testNote := notes.Note{Title: "Test", Status: notes.New}

		mockRepo.On("AddNote", mock.AnythingOfType("notes.Note")).Return(nil)

		id, err := service.CreateNote(testNote)

		require.NoError(t, err)
		assert.NotEmpty(t, id)
		_, uuidErr := uuid.Parse(id)
		assert.NoError(t, uuidErr)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)
		testNote := notes.Note{Title: "Test", Status: notes.New}

		mockRepo.On("AddNote", mock.AnythingOfType("notes.Note")).Return(errors.New("db error"))

		_, err := service.CreateNote(testNote)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNoteService_GetNotes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)
		expectedNotes := []notes.Note{
			{NID: "1", Title: "Note 1", Status: notes.New},
			{NID: "2", Title: "Note 2", Status: notes.Active},
		}

		mockRepo.On("GetNotes").Return(expectedNotes, nil)

		result, err := service.GetNotes()

		require.NoError(t, err)
		assert.Equal(t, expectedNotes, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)

		mockRepo.On("GetNotes").Return([]notes.Note{}, errors.New("db error"))

		_, err := service.GetNotes()

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNoteService_GetNoteID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)
		expectedNote := notes.Note{NID: "123", Title: "Test Note", Status: notes.Active}

		mockRepo.On("GetNoteID", "123").Return(expectedNote, nil)

		result, err := service.GetNoteID("123")

		require.NoError(t, err)
		assert.Equal(t, expectedNote, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)

		mockRepo.On("GetNoteID", "456").Return(notes.Note{}, errors.New("not found"))

		_, err := service.GetNoteID("456")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNoteService_DeleteNoteID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)

		mockRepo.On("DeleteNote", "123").Return(nil)

		err := service.DeleteNoteID("123")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)

		mockRepo.On("DeleteNote", "456").Return(errors.New("db error"))

		err := service.DeleteNoteID("456")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestNoteService_UpdateNoteID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)
		testNote := notes.Note{NID: "123", Title: "Updated", Status: notes.Active}

		mockRepo.On("UpdateNote", "123", testNote).Return(nil)

		err := service.UpdateNoteID("123", testNote)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewRepositoryNote(t)
		service := New(mockRepo)
		testNote := notes.Note{NID: "123", Title: "Updated", Status: notes.Active}

		mockRepo.On("UpdateNote", "123", testNote).Return(errors.New("db error"))

		err := service.UpdateNoteID("123", testNote)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
