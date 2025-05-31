package services

import (
	"encoding/json"
	"os"

	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"gorm.io/gorm"
)

type SyncDB struct {
	db        *gorm.DB
	notesFile string
}

func New(db *gorm.DB, filePath string) *SyncDB {
	return &SyncDB{
		db:        db,
		notesFile: filePath,
	}
}

func (s *SyncDB) SyncToDB() error {
	data, err := os.ReadFile(s.notesFile)
	if err != nil {
		return err
	}

	var fileNotes map[string]notes.Note
	if err := json.Unmarshal(data, &fileNotes); err != nil {
		return err
	}

	for _, note := range fileNotes {
		if err := s.db.Create(&note).Error; err != nil {
			continue
		}
	}
	return os.WriteFile(s.notesFile, []byte("{}"), 0644)
}
