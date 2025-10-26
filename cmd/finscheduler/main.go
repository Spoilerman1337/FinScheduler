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
	"net/http"
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

	repository := items.NewItemsRepository(db)
	service := items.NewItemsService(repository)
	handler := items.NewItemsHandler(service)

	r := chi.NewRouter()
	r.Mount("/api/items", handler.RegisterEndpoints())

	log.Println("Listening on :8080")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
