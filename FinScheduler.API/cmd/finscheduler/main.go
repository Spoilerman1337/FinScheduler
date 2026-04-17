package main

import (
	"context"
	"finscheduler/database"
	featurehttp "finscheduler/internal/features/http"
	"finscheduler/internal/features/repositories"
	"finscheduler/internal/features/services"
	"finscheduler/internal/health"
	"finscheduler/internal/infra"
	"finscheduler/internal/metrics"
	"finscheduler/internal/profiles"
	"finscheduler/internal/traces"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()
	cfg, err := infra.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	connectionString := cfg.ConnectionString

	mp, err := metrics.InitMetrics(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = mp.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown meter: %v", err)
		}
	}()
	metrics.InitInstruments()

	tp, err := traces.InitTracer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = tp.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown meter: %v", err)
		}
	}()

	prof, err := profiles.InitProfiler()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = prof.Stop()
		if err != nil {
			log.Fatalf("failed to shutdown profiler: %v", err)
		}
	}()

	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.RunMigrations(connectionString)

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(stdoutHandler)

	itemsRepository := repositories.NewItemsRepository(db, logger)
	tagsRepository := repositories.NewTagsRepository(db, logger)
	tagsToItemsRepository := repositories.NewTagToItemsRepository(db, logger)

	itemsService := services.NewItemsService(itemsRepository, tagsRepository, tagsToItemsRepository, logger)
	tagsService := services.NewTagsService(tagsRepository, logger)

	tagsHandler := featurehttp.NewTagsHandler(tagsService, logger)
	itemsHandler := featurehttp.NewItemsHandler(itemsService, logger)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))
	health.SetupHealthChecks(r, db)
	r.Route("/api/items", func(r chi.Router) {
		itemsHandler.RegisterEndpoints(r)
	})
	r.Route("/api/tags", func(r chi.Router) {
		tagsHandler.RegisterEndpoints(r)
	})

	log.Printf("Listening at :%d", cfg.ServerPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
