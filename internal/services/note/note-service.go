package note

import (
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"

	"github.com/google/uuid"
)

type RepositoryNote interface {
	AddNote(note notes.Note) error
	GetNotes() (map[string]notes.Note, error)
	GetNoteID(noteID string) (notes.Note, error)
	DeleteNote(noteID string) error
	UpdateNote(noteID string, note notes.Note) error
}

type NoteService struct {
	repo RepositoryNote
}

func New(repo RepositoryNote) *NoteService {
	return &NoteService{repo: repo}
}
func (ns *NoteService) CreateNote(note notes.Note) (string, error) {
	note.NID = uuid.New().String()

	err := ns.repo.AddNote(note)
	if err != nil {
		return ``, err
	}
	return note.NID, nil
}

func (ns *NoteService) GetNotes() (map[string]notes.Note, error) {
	notes, err := ns.repo.GetNotes()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (ns *NoteService) GetNoteID(noteID string) (notes.Note, error) {
	note, err := ns.repo.GetNoteID(noteID)
	if err != nil {
		return notes.Note{}, err
	}
	return note, nil
}

func (ns *NoteService) DeleteNoteID(noteID string) error {
	err := ns.repo.DeleteNote(noteID)
	if err != nil {
		return err
	}
	return nil
}

func (ns *NoteService) UpdateNoteID(noteID string, note notes.Note) error {
	err := ns.repo.UpdateNote(noteID, note)
	if err != nil {
		return err
	}
	return nil
}