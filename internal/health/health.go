package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SetupHealthChecks(db *sqlx.DB) {
	router := chi.NewRouter()
	router.Handle("/livez", LiveHandler())
	router.Handle("/readyz", ReadyHandler(db))
}
