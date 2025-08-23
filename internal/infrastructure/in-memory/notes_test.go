package inmemory

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noteComparer struct {
	NID         string
	Title       string
	Description string
	Status      notes.Status
	UID         string
	Deleted     bool
}

func toComparer(n notes.Note) noteComparer {
	return noteComparer{
		NID:         n.NID,
		Title:       n.Title,
		Description: n.Description,
		Status:      n.Status,
		UID:         n.UID,
		Deleted:     n.Deleted,
	}
}

func TestAddNote(t *testing.T) {
	tmpFile := t.TempDir() + "/notes_test.json"
	im := NewNotes(false, tmpFile)

	testNote := notes.Note{
		NID:         "1",
		Title:       "Test Note",
		Description: "Test Content",
		Status:      notes.New,
		CreatedAt:   time.Now(),
		UID:         "user1",
	}

	t.Run("successful add", func(t *testing.T) {
		err := im.AddNote(testNote)
		assert.NoError(t, err)
		assert.Contains(t, im.noteStorage, testNote.NID)
		assert.Equal(t, toComparer(testNote), toComparer(im.noteStorage[testNote.NID]))

		data, err := os.ReadFile(tmpFile)
		assert.NoError(t, err)

		var fileData map[string]notes.Note
		err = json.Unmarshal(data, &fileData)
		assert.NoError(t, err)
		assert.Equal(t, toComparer(testNote), toComparer(fileData[testNote.NID]))
	})

	t.Run("duplicate title", func(t *testing.T) {
		duplicateNote := testNote
		duplicateNote.NID = "2"
		err := im.AddNote(duplicateNote)
		assert.ErrorIs(t, err, notes.ErrNoteAlreadyExists)
	})
}

func TestGetNotes(t *testing.T) {
	tmpFile := t.TempDir() + "/notes_test.json"
	im := NewNotes(false, tmpFile)

	now := time.Now()
	testNotes := []notes.Note{
		{
			NID:         "1",
			Title:       "Note 1",
			Description: "Content 1",
			Status:      notes.Active,
			CreatedAt:   now,
			UID:         "user1",
		},
		{
			NID:         "2",
			Title:       "Note 2",
			Description: "Content 2",
			Status:      notes.Inactive,
			CreatedAt:   now,
			UID:         "user1",
		},
	}

	for _, note := range testNotes {
		err := im.AddNote(note)
		require.NoError(t, err)
	}

	t.Run("get all notes", func(t *testing.T) {
		result, err := im.GetNotes()
		assert.NoError(t, err)
		assert.ElementsMatch(t, testNotes, result)
	})

	t.Run("empty storage", func(t *testing.T) {
		emptyIm := NewNotes(false, t.TempDir()+"/empty.json")
		result, err := emptyIm.GetNotes()
		assert.ErrorIs(t, err, notes.ErrNoNotesAvailable)
		assert.Nil(t, result)
	})
}

func TestGetNoteID(t *testing.T) {
	tmpFile := t.TempDir() + "/notes_test.json"
	im := NewNotes(false, tmpFile)

	now := time.Now()
	testNote := notes.Note{
		NID:         "1",
		Title:       "Test Note",
		Description: "Test Content",
		Status:      notes.Active,
		CreatedAt:   now,
		UID:         "user1",
	}

	err := im.AddNote(testNote)
	require.NoError(t, err)

	t.Run("note exists", func(t *testing.T) {
		result, getErr := im.GetNoteID("1")
		assert.NoError(t, getErr)
		assert.Equal(t, testNote, result)
	})

	t.Run("note not found", func(t *testing.T) {
		result, getErr := im.GetNoteID("999")
		assert.ErrorIs(t, getErr, notes.ErrNoteNotFound)
		assert.Equal(t, notes.Note{}, result)
	})
}

func TestDeleteNote(t *testing.T) {
	tmpFile := t.TempDir() + "/notes_test.json"
	im := NewNotes(false, tmpFile)

	now := time.Now()
	testNote := notes.Note{
		NID:         "1",
		Title:       "Test Note",
		Description: "Test Content",
		Status:      notes.Active,
		CreatedAt:   now,
		UID:         "user1",
	}

	err := im.AddNote(testNote)
	require.NoError(t, err)

	t.Run("successful delete", func(t *testing.T) {
		err = im.DeleteNote("1")
		assert.NoError(t, err)
		assert.NotContains(t, im.noteStorage, "1")

		data, readErr := os.ReadFile(tmpFile)
		assert.NoError(t, readErr)

		var fileData map[string]notes.Note
		err = json.Unmarshal(data, &fileData)
		assert.NoError(t, err)
		assert.NotContains(t, fileData, "1")
	})

	t.Run("note not found", func(t *testing.T) {
		err = im.DeleteNote("999")
		assert.ErrorIs(t, err, notes.ErrNoteNotFound)
	})
}

func TestUpdateNote(t *testing.T) {
	tmpFile := t.TempDir() + "/notes_test.json"
	im := NewNotes(false, tmpFile)

	initialNote := notes.Note{
		NID:         "1",
		Title:       "Old Title",
		Description: "Old Content",
		Status:      notes.Active,
		CreatedAt:   time.Now(),
		UID:         "user1",
	}

	err := im.AddNote(initialNote)
	require.NoError(t, err)

	updatedNote := initialNote
	updatedNote.Title = "New Title"
	updatedNote.Description = "New Content"
	updatedNote.Status = notes.Inactive

	t.Run("successful update", func(t *testing.T) {
		err = im.UpdateNote("1", updatedNote)
		assert.NoError(t, err)
		assert.Equal(t, toComparer(updatedNote), toComparer(im.noteStorage["1"]))

		data, readErr := os.ReadFile(tmpFile)
		assert.NoError(t, readErr)

		var fileData map[string]notes.Note
		err = json.Unmarshal(data, &fileData)
		assert.NoError(t, err)
		assert.Equal(t, toComparer(updatedNote), toComparer(fileData["1"]))
	})

	t.Run("note not found", func(t *testing.T) {
		err = im.UpdateNote("999", updatedNote)
		assert.ErrorIs(t, err, notes.ErrNoteNotFound)
	})
}

func TestLoadFromFile(t *testing.T) {
	tmpFile := t.TempDir() + "/notes_load_test.json"

	testNote := notes.Note{
		NID:         "1",
		Title:       "Loaded Note",
		Description: "Loaded Content",
		Status:      notes.New,
		CreatedAt:   time.Now(),
		UID:         "user1",
	}

	data := map[string]notes.Note{"1": testNote}
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	err = os.WriteFile(tmpFile, jsonData, 0644)
	require.NoError(t, err)

	im := NewNotes(false, tmpFile)

	t.Run("load from existing file", func(t *testing.T) {
		assert.Contains(t, im.noteStorage, "1")
		assert.Equal(t, toComparer(testNote), toComparer(im.noteStorage["1"]))
	})

	t.Run("load from non-existent file", func(t *testing.T) {
		newIm := NewNotes(false, "nonexistent.json")
		assert.Empty(t, newIm.noteStorage)
	})
}

func TestStatusOperations(t *testing.T) {
	tests := []struct {
		name           string
		parseInput     string
		expectedStatus notes.Status
	}{
		{"Status New", "New", notes.New},
		{"Status Active", "Active", notes.Active},
		{"Status Inactive", "Inactive", notes.Inactive},
		{"Status Deleted", "Deleted", notes.Deleted},
		{"Invalid Status", "Invalid", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := notes.ParseStatus(tt.parseInput)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}
