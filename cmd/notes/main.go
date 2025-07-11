package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Snoop-Duck/ToDoList/internal/server"
	"github.com/Snoop-Duck/ToDoList/internal/services"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	dbstorage "github.com/Snoop-Duck/ToDoList/internal/infrastructure/db-storage"
	inmemory "github.com/Snoop-Duck/ToDoList/internal/infrastructure/in-memory"

	"github.com/Snoop-Duck/ToDoList/internal"
	logger "github.com/Snoop-Duck/ToDoList/pkg"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
)

func gracefulShutdown(cancel context.CancelFunc) {
	log := logger.Get()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	sig := <-c

	log.Info().Msgf("graceful shutdown with signal: %s", sig)

	cancel()
}

func main() {
	cfg, err := internal.ReadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.Get(cfg.Debug)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gracefulShutdown(cancel)

	log.Info().Msg("service starting")

	dns := os.Getenv("DB_CONNECTION_STRING")
	if dns == "" {
		dns = "postgres://user:password@db:5432/notes?sslmode=disable"
	}
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
	if err = dbstorage.ApplyMigrations(dns); err != nil {
		log.Warn().Err(err).Msg("failed to apply migrations. Use in memory storage")
		repoUser.Close()
		repoUser = inmemory.NewUsers()
	}

	syncService := services.New(db, "storage/notes.json")

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := syncService.SyncToDB(); err != nil {
					log.Error().Err(err).Msg("Sync failed")
				} else {
					log.Info().Msg("Sync completed successfully")
				}
			case <-ctx.Done():
				log.Info().Msg("Stopping sync service")
				return
			}
		}
	}()

	notesAPI := server.New(cfg, repoUser, repoNote)

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		log.Info().Str("address", cfg.Host+":"+strconv.Itoa(cfg.Port)).Msg("server starting")
		if err := notesAPI.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := notesAPI.Stop(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("failed to stop server gracefully")
			return err
		}
		if err := repoUser.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close user repository")
			return err
		}

		if err := syncService.SyncToDB(); err != nil {
			log.Error().Err(err).Msg("final sync failed")
		}

		return nil
	})

	if err := group.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("service stopped with error")
		}
	}
	log.Info().Msg("service stopped gracefully")
}
