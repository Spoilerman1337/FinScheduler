package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SetupHealthChecks(router *chi.Mux, db *sqlx.DB) {
	router.Handle("/livez", LiveHandler())
	router.Handle("/readyz", ReadyHandler(db))
}
