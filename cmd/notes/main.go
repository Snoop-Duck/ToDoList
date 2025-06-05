package main

import (
	"context"
	"strconv"
	"time"

	"github.com/Snoop-Duck/ToDoList/internal/server"
	"github.com/Snoop-Duck/ToDoList/internal/services"
	"gorm.io/gorm"

	dbstorage "github.com/Snoop-Duck/ToDoList/internal/infrastructure/db-storage"
	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/in-memory"

	"github.com/Snoop-Duck/ToDoList/internal"
	logger "github.com/Snoop-Duck/ToDoList/pkg"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
)

func main() {
	cfg := internal.ReadConfig()

	log := logger.Get(cfg.Debug)

	log.Info().Msg("service starting")

	dns := "postgres://user:password@localhost:5432/notes?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	var repoUser server.Repository
	repoNote := inmemory.NewNotes(cfg.Debug, "storage/notes.json")
	repoUser, err = dbstorage.New(context.Background(), dns)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to db. Use in memory storage")
		repoUser = inmemory.NewUsers()
	}
	if err = dbstorage.AppyMigrations(dns); err != nil {
		log.Warn().Err(err).Msg("failed to apply migrations. Use in memory storage")
		repoUser.Close()
		repoUser = inmemory.NewUsers()
	}

	syncService := services.New(db, "storage/notes.json")

	go func() {
		for range time.Tick(5 * time.Minute) {
			if err := syncService.SyncToDB(); err != nil {
				log.Error().Err(err).Msg("Sync failed")
			} else {
				log.Info().Msg("Sync completed successfully")
			}
		}
	}()

	notesAPI := server.New(cfg, repoUser, repoNote)
	log.Info().Str("address", cfg.Host+":"+strconv.Itoa(cfg.Port)).Msg("server started")
	if err := notesAPI.Run(); err != nil {
		log.Error().Err(err).Msg("fatal running server")
	}
}
