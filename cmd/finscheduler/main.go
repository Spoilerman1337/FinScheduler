package main

import (
	"finscheduler/database"
	"finscheduler/internal/infra"
	"finscheduler/internal/items"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg, err := infra.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	connectionString := cfg.ConnectionString

	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database.RunMigrations(connectionString)

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(stdoutHandler)

	repository := items.NewItemsRepository(db, logger)
	service := items.NewItemsService(repository, logger)
	handler := items.NewItemsHandler(service, logger)

	r := chi.NewRouter()
	r.Mount("/api/items", handler.RegisterEndpoints())

	log.Println("Listening on :8080")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
