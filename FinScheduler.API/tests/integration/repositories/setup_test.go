//go:build integration
// +build integration

package repositories_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"finscheduler/tests/internal/testsupport"

	"github.com/jmoiron/sqlx"
)

var testDB *sqlx.DB
var testLogger *slog.Logger
var testContext context.Context

const closedDBDriverName = "pgx"
const closedDBConnectionString = "postgres://test:secret@127.0.0.1:1/testdb?sslmode=disable&connect_timeout=1"

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

func newClosedDB(t testing.TB) *sqlx.DB {
	t.Helper()

	db := sqlx.MustOpen(closedDBDriverName, closedDBConnectionString)
	closeErr := db.Close()
	requireNoError(t, closeErr)

	return db
}

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
