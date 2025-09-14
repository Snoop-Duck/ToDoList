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

	_ "github.com/Snoop-Duck/ToDoList/docs"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	contextTimeout = 5 * time.Second
)

func gracefulShutdown(cancel context.CancelFunc) {
	log := logger.Get()

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	sig := <-c

	log.Info().Msgf("graceful shutdown with signal: %s", sig)

	cancel()
}

func setupDatabase(log logger.Logger, dns string, debug bool) (server.Repository, server.RepositoryNote, *gorm.DB, error) {
	ctx := context.Background()

	// 1. Пытаемся подключиться к БД через pgxpool
	dbPool, err := pgxpool.New(ctx, dns)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to database. Using in-memory storage")
		return inmemory.NewUsers(debug), inmemory.NewNotes(debug, "storage/notes.json"), nil, nil
	}
	defer dbPool.Close()

	// 2. Пытаемся применить миграции
	if err = dbstorage.ApplyMigrations(dns); err != nil {
		log.Warn().Err(err).Msg("failed to apply migrations. Using in-memory storage")
		return inmemory.NewUsers(debug), inmemory.NewNotes(debug, "storage/notes.json"), nil, nil
	}

	// 3. Создаем DBStorage который будет работать с БД
	dbStorage, err := dbstorage.New(ctx, dns)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create DB storage. Using in-memory storage")
		return inmemory.NewUsers(debug), inmemory.NewNotes(debug, "storage/notes.json"), nil, nil
	}

	// 4. Также создаем gorm connection для sync service
	gormDB, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Warn().Err(err).Msg("failed to create gorm connection, sync service will not work")
		// Но продолжаем работу с БД репозиториями
		return dbStorage, dbStorage, nil, nil
	}

	// 5. Возвращаем ОБА репозитория из DBStorage
	return dbStorage, dbStorage, gormDB, nil
}

func startSyncService(ctx context.Context, db *gorm.DB, log logger.Logger) {
	syncService := services.New(db, "storage/notes.json")

	go func() {
		ticker := time.NewTicker(contextTimeout)
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
}

func runServer(
	ctx context.Context,
	cfg *internal.Config,
	notesAPI *server.NotesAPI,
	repoUser server.Repository,
	db *gorm.DB,
	log logger.Logger,
) error {
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

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), contextTimeout)
		defer shutdownCancel()

		if err := notesAPI.Stop(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("failed to stop server gracefully")
		}

		if err := repoUser.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close user repository")
		}

		if db != nil {
			syncService := services.New(db, "storage/notes.json")
			if err := syncService.SyncToDB(); err != nil {
				log.Error().Err(err).Msg("final sync failed")
			}
		}

		return nil
	})

	return group.Wait()
}

// @title ToDoList API
// @version 1.0
// @description API для управления задачами и заметками
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@todolist.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT token for authentication. Format: "Bearer {token}"
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
		dns = "postgres://user:password@localhost:5432/notes?sslmode=disable"
	}

	// Получаем оба репозитория и gorm соединение
	repoUser, repoNote, db, setupErr := setupDatabase(log, dns, cfg.Debug)
	if setupErr != nil {
		log.Error().Err(setupErr).Msg("failed to setup database")
		return
	}

	// Запускаем sync service только если есть gorm соединение
	if db != nil {
		startSyncService(ctx, db, log)
	}

	notesAPI := server.New(cfg, repoUser, repoNote)

	if runErr := runServer(ctx, cfg, notesAPI, repoUser, db, log); runErr != nil {
		if !errors.Is(runErr, http.ErrServerClosed) {
			log.Error().Err(runErr).Msg("service stopped with error")
			return
		}
	}

	log.Info().Msg("service stopped gracefully")
}
