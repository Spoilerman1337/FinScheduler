//go:build integration
// +build integration

package services_test

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
