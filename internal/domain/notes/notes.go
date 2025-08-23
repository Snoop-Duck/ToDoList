package notes

import "time"

type Note struct {
	NID         string    `json:"nid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UID         string    `json:"uid"`
	Deleted     bool      `json:"deleted"`
}
type NoteResponseFormat struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UID         string `json:"uid"`
}

func NoteResponse(note Note) NoteResponseFormat {
	return NoteResponseFormat{
		Title:       note.Title,
		Description: note.Description,
		Status:      note.Status.String(),
		CreatedAt:   note.CreatedAt.Format(time.RFC3339),
		UID:         note.UID,
	}
}
