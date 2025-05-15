package inmemory

import (
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
)

type InMemoryNotes struct {
	noteStorage map[string]notes.Note
}

func New() *InMemoryNotes {
	return &InMemoryNotes{
		noteStorage: make(map[string]notes.Note),
	}
}

func (im *InMemoryNotes) AddNote(note notes.Note) error {
	for _, id := range im.noteStorage {
		if id.Title == note.Title {
			return notes.ErrNoteAlredyExists
		}
	}
	im.noteStorage[note.UID] = note
	return nil
}

func (im *InMemoryNotes) GetNotes() (map[string]notes.Note, error) {
	if len(im.noteStorage) == 0 {
		return nil, notes.ErrNoNotesAvailable
	}
	return im.noteStorage, nil
}

func (im *InMemoryNotes) GetNoteID(noteID string) (notes.Note, error) {
	note, ok := im.noteStorage[noteID]
	if !ok {
		return notes.Note{}, notes.ErrNoteNotFound
	}
	return note, nil
}

func (im *InMemoryNotes) DeleteNote(noteID string) error {
	if _, ok := im.noteStorage[noteID]; !ok {
		return notes.ErrNoteNotFound
	}
	delete(im.noteStorage, noteID)
	return nil
}

func (im *InMemoryNotes) UpdateNote(noteID string, note notes.Note) error {
	if _, ok := im.noteStorage[noteID]; !ok {
		return notes.ErrNoteNotFound
	}
	im.noteStorage[noteID] = note
	return nil
}
