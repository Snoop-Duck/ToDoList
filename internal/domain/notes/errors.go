package notes

import "errors"

var (
	ErrNoteNotFound     = errors.New("note not found")
	ErrNoNotesAvailable = errors.New("no notes available")
	ErrNoteAlredyExists = errors.New("note alredy exists")
)
