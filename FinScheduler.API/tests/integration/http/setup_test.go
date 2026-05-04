//go:build integration
// +build integration

package featurehttp_test

import (
	"context"
	featurehttp "finscheduler/internal/features/http"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

var testDB *sqlx.DB
var testLogger *slog.Logger
var testContext context.Context

type testApplication struct {
	router       http.Handler
	itemsService *services.ItemsService
	tagsService  *services.TagsService
}

func TestMain(m *testing.M) {
	env, err := testsupport.NewPostgresEnvironment(context.Background())
	if err != nil {
		fmt.Printf("failed to start postgres: %v", err)
		os.Exit(1)
	}

	testDB = env.DB
	testLogger = env.Logger
	testContext = env.Context

	code := m.Run()

	if err := env.Close(); err != nil {
		fmt.Printf("failed to stop postgres: %v", err)
	}

	os.Exit(code)
}

func newTestApplication() *testApplication {
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	tagsService := services.NewTagsService(uow, testLogger)
	itemsHandler := featurehttp.NewItemsHandler(itemsService, testLogger)
	tagsHandler := featurehttp.NewTagsHandler(tagsService, testLogger)
	router := chi.NewRouter()

	router.Route("/api/items", func(route chi.Router) {
		itemsHandler.RegisterEndpoints(route)
	})
	router.Route("/api/tags", func(route chi.Router) {
		tagsHandler.RegisterEndpoints(route)
	})

	return &testApplication{
		router:       router,
		itemsService: itemsService,
		tagsService:  tagsService,
	}
}

func newJSONRequest(method string, target string, body string) *http.Request {
	bodyReader := strings.NewReader(body)
	request := httptest.NewRequest(method, target, bodyReader)

	if len(body) > 0 {
		contentTypeHeader := "Content-Type"
		contentTypeValue := "application/json"
		request.Header.Set(contentTypeHeader, contentTypeValue)
	}

	return request
}
