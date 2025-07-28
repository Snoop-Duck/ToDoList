package notes

import "errors"

var (
	ErrNoteNotFound      = errors.New("note not found")
	ErrNoNotesAvailable  = errors.New("no notes available")
	ErrNoteAlreadyExists = errors.New("note already exists")
)
