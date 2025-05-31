package inmemory

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Dorrrke/notes-g2/pkg/logger"
	"github.com/Snoop-Duck/ToDoList/internal/domain/notes"
	"github.com/Snoop-Duck/ToDoList/internal/domain/users"
	"github.com/rs/zerolog"
)

type InMemoryNotes struct {
	noteStorage map[string]notes.Note
	filePath    string
	log         zerolog.Logger
}

var emtyUser = users.User{}

type InMemoryUsers struct {
	userStorage map[string]users.User
	log         zerolog.Logger
}

func NewNotes(filePath string) *InMemoryNotes {
	storage := &InMemoryNotes{
		noteStorage: make(map[string]notes.Note),
		filePath:    filePath,
		log:         logger.Get(false),
	}

	if err := storage.loadFromFile(); err != nil {
		storage.log.Error().Err(err).Msg("Не удалось загрузить данные из файла")
	}
	return storage
}

func NewUsers() *InMemoryUsers {
	return &InMemoryUsers{
		userStorage: make(map[string]users.User),
		log:         logger.Get(false),
	}
}

func (im *InMemoryNotes) loadFromFile() error {
	if _, err := os.Stat(im.filePath); os.IsNotExist(err) {
		im.log.Debug().Msg("Файл хранилища не существует, будет создан новый")
		return nil
	}

	data, err := os.ReadFile(im.filePath)
	if err != nil {
		im.log.Error().Err(err).Msg("Ошибка чтения файла")
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if err := json.Unmarshal(data, &im.noteStorage); err != nil {
		im.log.Error().Err(err).Msg("Ошибка парсинга JSON")
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	im.log.Info().Int("count", len(im.noteStorage)).Msg("Заметки успешно загружены из файла")
	return nil
}

func (im *InMemoryNotes) SaveToFile() error {
	data, err := json.MarshalIndent(im.noteStorage, "", " ")
	if err != nil {
		im.log.Error().Err(err).Msg("Ошибка сериализации")
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	if err := os.WriteFile(im.filePath, data, 0644); err != nil {
		im.log.Error().Err(err).Msg("Ошибка записи в файл")
		return fmt.Errorf("ошибка записи: %w", err)
	}
	im.log.Info().Int("count", len(im.noteStorage)).Msg("Заметки успешно сохранены в файл")
	return nil
}
