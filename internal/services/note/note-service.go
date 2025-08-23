package note

import (
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"

	"github.com/google/uuid"
)

type RepositoryNote interface {
	AddNote(note notes.Note) error
	GetNotes() ([]notes.Note, error)
	GetNoteID(noteID string) (notes.Note, error)
	DeleteNote(noteID string) error
	UpdateNote(noteID string, note notes.Note) error
}

type Service struct {
	repo RepositoryNote
}

func New(repo RepositoryNote) *Service {
	return &Service{repo: repo}
}
func (ns *Service) CreateNote(note notes.Note) (string, error) {
	note.NID = uuid.New().String()

	err := ns.repo.AddNote(note)
	if err != nil {
		return ``, err
	}
	return note.NID, nil
}

func (ns *Service) GetNotes() ([]notes.Note, error) {
	notes, err := ns.repo.GetNotes()
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (ns *Service) GetNoteID(noteID string) (notes.Note, error) {
	note, err := ns.repo.GetNoteID(noteID)
	if err != nil {
		return notes.Note{}, err
	}
	return note, nil
}

func (ns *Service) DeleteNoteID(noteID string) error {
	err := ns.repo.DeleteNote(noteID)
	if err != nil {
		return err
	}
	return nil
}

func (ns *Service) UpdateNoteID(noteID string, note notes.Note) error {
	err := ns.repo.UpdateNote(noteID, note)
	if err != nil {
		return err
	}
	return nil
}
