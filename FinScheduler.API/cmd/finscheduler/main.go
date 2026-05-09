package main

import (
	"context"
	"finscheduler/database"
	featurehttp "finscheduler/internal/features/http"
	"finscheduler/internal/features/services"
	"finscheduler/internal/health"
	"finscheduler/internal/infra"
	"finscheduler/internal/logging"
	"finscheduler/internal/metrics"
	"finscheduler/internal/persistence"
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

	tp, err := traces.InitTracer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = tp.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to shutdown meter: %v", err)
		}
	}()

	prof, err := profiles.InitProfiler(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if prof != nil {
		defer func() {
			err = prof.Stop()
			if err != nil {
				log.Fatalf("failed to shutdown profiler: %v", err)
			}
		}()
	}

	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.RunMigrations(connectionString)

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(logging.NewCustomLoggingHandler(stdoutHandler))

	uow := persistence.NewUnitOfWork(db, logger)

	itemsService := services.NewItemsService(uow, logger)
	tagsService := services.NewTagsService(uow, logger)

	tagsHandler := featurehttp.NewTagsHandler(tagsService, logger)
	itemsHandler := featurehttp.NewItemsHandler(itemsService, logger)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSSettings.AllowedOrigins,
		AllowedMethods:   cfg.CORSSettings.AllowedMethods,
		AllowedHeaders:   cfg.CORSSettings.AllowedHeaders,
		AllowCredentials: cfg.CORSSettings.AllowCredentials,
	}))
	if cfg.Observability.Metrics.Enabled {
		r.Handle(cfg.Observability.Metrics.ExportEndpoint, metrics.Handler())
	}
	health.SetupHealthChecks(r, db)
	r.Route("/api/items", func(r chi.Router) {
		itemsHandler.RegisterEndpoints(r)
	})
	r.Route("/api/tags", func(r chi.Router) {
		tagsHandler.RegisterEndpoints(r)
	})

	logger.Info("starting http server",
		"port", cfg.ServerPort,
		"observability_service_name", cfg.Observability.ServiceName,
		"metrics_enabled", cfg.Observability.Metrics.Enabled,
		"metrics_export_endpoint", cfg.Observability.Metrics.ExportEndpoint,
		"traces_enabled", cfg.Observability.Traces.Enabled,
		"trace_export_endpoint", cfg.Observability.Traces.ExportEndpoint,
		"profiling_enabled", cfg.Observability.Profiling.Enabled,
	)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
