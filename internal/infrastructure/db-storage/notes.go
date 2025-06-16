package dbstorage

import (
	"context"
	"time"

	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
)

func (db *DBStorage) AddNote(notes notes.Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "INSERT INTO note(nid, title, description, status, created_at, user_id)")
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStorage) GetNotes() ([]notes.Note, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.db.Query(ctx, "SELECT * FROM notes")
	if err != nil {
		log.Error().Err(err).Msg("failed to get notes")
		return nil, err
	}

	var notesSlice []notes.Note
	for rows.Next() {
		var note notes.Note
		if err := rows.Scan(&note.NID, &note.Title, &note.Description, &note.Status, &note.Created_at, &note.UID); err != nil {
			log.Error().Err(err).Msg("failed to scan note")
			return nil, err
		}
		notesSlice = append(notesSlice, note)
	}
	return notesSlice, nil
}

func (db *DBStorage) GetNoteID(noteID string) (notes.Note, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var note notes.Note

	row := db.db.QueryRow(ctx, "SELECT * FROM notes WHERE nid = $1", noteID)
	err := row.Scan(&note.NID, &note.Title, &note.Description, &note.Status, &note.Created_at, &note.UID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get note")
		return notes.Note{}, err
	}
	return note, nil
}

func (db *DBStorage) DeleteNote(noteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "DELETE FROM notes WHERE nid = $1", noteID)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStorage) UpdateNote(noteID string, note notes.Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "UPDATE notes SET title = $1, description = $2, status = $3 WHERE nid = $4", note.Title, note.Description, note.Status, note.NID)
	if err != nil {
		return err
	}
	return nil
}
