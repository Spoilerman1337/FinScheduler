//go:build integration
// +build integration

package httpapi_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	featurehttp "finscheduler/internal/features/http"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

var testDB *sqlx.DB
var testLogger *slog.Logger
var testContext context.Context
var testFixtures testsupport.Fixtures

func TestMain(m *testing.M) {
	env, err := testsupport.NewPostgresEnvironment(context.Background())
	if err != nil {
		fmt.Printf("failed to start postgres: %v", err)
		os.Exit(1)
	}

	testDB = env.DB
	testLogger = env.Logger
	testContext = env.Context
	testFixtures = testsupport.NewFixtures(testContext, testDB, testLogger)

	code := m.Run()

	if err := env.Close(); err != nil {
		fmt.Printf("failed to stop postgres: %v", err)
	}

	os.Exit(code)
}

func newTestRouter() *chi.Mux {
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsHandler := featurehttp.NewItemsHandler(services.NewItemsService(uow, testLogger), testLogger)
	tagsHandler := featurehttp.NewTagsHandler(services.NewTagsService(uow, testLogger), testLogger)

	router := chi.NewRouter()
	router.Route("/api/items", func(r chi.Router) {
		itemsHandler.RegisterEndpoints(r)
	})
	router.Route("/api/tags", func(r chi.Router) {
		tagsHandler.RegisterEndpoints(r)
	})

	return router
}
