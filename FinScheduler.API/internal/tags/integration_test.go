//go:build integration
// +build integration

package tags

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"finscheduler/internal/infra"
	"finscheduler/internal/metrics"

	"log"
)

// TODO: Инициализация тестов достаточно одинаковая, нужно выделить общую логику, чтобы не поддерживать одно и то же
var testDB *sqlx.DB
var testLogger *slog.Logger
var testContext context.Context

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, terminate, db, err := setupPostgresContainer(ctx)
	if err != nil {
		fmt.Printf("failed to start postgres: %v", err)
	}

	mp, err := metrics.InitMetrics(ctx, &infra.Config{Env: "Testing", ServiceName: "fin-scheduler-api"})
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

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(stdoutHandler)

	testDB = db
	testLogger = logger
	testContext = ctx

	tagSchema := `
        CREATE TABLE tags(
			id        UUID PRIMARY KEY,
			name      TEXT    NOT NULL UNIQUE,
			is_active BOOLEAN NOT NULL DEFAULT FALSE
		);
    `

	setupSchema(tagSchema, "tags")

	code := m.Run()

	terminate(ctx, container)
	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, func(context.Context, testcontainers.Container), *sqlx.DB, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://test:secret@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, nil, nil, err
	}

	terminate := func(ctx context.Context, c testcontainers.Container) {
		_ = c.Terminate(ctx)
	}

	return container, terminate, db, nil
}

func setupSchema(sqlSchema string, name string) error {
	_, err := testDB.Exec(sqlSchema)
	if err != nil {
		fmt.Printf("failed to create %s schema: %v", name, err)
	}

	return err
}
