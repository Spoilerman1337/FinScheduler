//go:build integration
// +build integration

package testsupport

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"finscheduler/internal/infra"
	"finscheduler/internal/metrics"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Environment struct {
	Context context.Context
	DB      *sqlx.DB
	Logger  *slog.Logger

	container     testcontainers.Container
	meterProvider *sdkmetric.MeterProvider
}

func NewPostgresEnvironment(ctx context.Context) (*Environment, error) {
	container, db, err := setupPostgresContainer(ctx)
	if err != nil {
		return nil, err
	}

	mp, err := metrics.InitMetrics(ctx, &infra.Config{
		Env: "Testing",
		Observability: infra.ObservabilityConfig{
			ServiceName: "fin-scheduler-api",
		},
	})
	if err != nil {
		_ = db.Close()
		_ = container.Terminate(ctx)
		return nil, err
	}

	metrics.InitInstruments()

	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	env := &Environment{
		Context:       ctx,
		DB:            db,
		Logger:        slog.New(stdoutHandler),
		container:     container,
		meterProvider: mp,
	}

	if err := setupSchema(env.DB); err != nil {
		_ = env.Close()
		return nil, err
	}

	return env, nil
}

func (env *Environment) Close() error {
	var closeErr error
	if env.meterProvider != nil {
		if err := env.meterProvider.Shutdown(env.Context); err != nil {
			log.Printf("failed to shutdown meter: %v", err)
		}
	}
	if env.DB != nil {
		closeErr = env.DB.Close()
	}
	if env.container != nil {
		if err := env.container.Terminate(env.Context); err != nil && closeErr == nil {
			closeErr = err
		}
	}

	return closeErr
}

func Truncate(t testing.TB, db *sqlx.DB, tables ...string) {
	t.Helper()

	if len(tables) == 0 {
		tables = []string{"items", "tags", "tag_to_item"}
	}

	query := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, *sqlx.DB, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_USER":     "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		).WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	dsn := fmt.Sprintf("postgres://test:secret@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	if err := waitForDatabase(ctx, db); err != nil {
		_ = db.Close()
		_ = container.Terminate(ctx)
		return nil, nil, err
	}

	return container, db, nil
}

func waitForDatabase(ctx context.Context, db *sqlx.DB) error {
	deadline := time.Now().Add(30 * time.Second)
	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		} else if time.Now().After(deadline) {
			return fmt.Errorf("wait for postgres readiness: %w", err)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for postgres readiness: %w", ctx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func setupSchema(db *sqlx.DB) error {
	if err := setupItemsSchema(db); err != nil {
		return err
	}
	if err := setupPriceHistorySchema(db); err != nil {
		return err
	}
	if err := setupPriceForecastSchema(db); err != nil {
		return err
	}
	if err := setupTagsSchema(db); err != nil {
		return err
	}
	if err := setupTagToItemSchema(db); err != nil {
		return err
	}

	return nil
}

func setupItemsSchema(db *sqlx.DB) error {
	return setupTable(db, "items", `
		CREATE TABLE items (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			price NUMERIC(16, 2) NOT NULL DEFAULT 0 CHECK (price >= 0),
			description TEXT NULL,
			is_active BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NULL,
			cashback INTEGER NOT NULL DEFAULT 0,
			category TEXT NOT NULL DEFAULT 'None'
		);
	`)
}

func setupPriceHistorySchema(db *sqlx.DB) error {
	return setupTable(db, "price_history", `
		CREATE TABLE price_history (
			id UUID PRIMARY KEY,
			item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
			recorded_at DATE NOT NULL,
			value NUMERIC(16, 2) NOT NULL CHECK (value >= 0),
			CONSTRAINT uq_price_history_item_id_recorded_at
				UNIQUE (item_id, recorded_at)
		);

		CREATE INDEX idx_price_history_item_id
			ON price_history (item_id);
	`)
}

func setupPriceForecastSchema(db *sqlx.DB) error {
	return setupTable(db, "price_forecast", `
		CREATE TABLE price_forecast (
			id UUID PRIMARY KEY,
			item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
			calculated_at DATE NOT NULL,
			last_known_price NUMERIC(16, 2) NOT NULL CHECK (last_known_price >= 0),
			average_monthly_drift NUMERIC(16, 2) NOT NULL,
			CONSTRAINT uq_price_forecast_item_id_calculated_at
				UNIQUE (item_id, calculated_at)
		);
	`)
}

func setupTagsSchema(db *sqlx.DB) error {
	return setupTable(db, "tags", `
		CREATE TABLE tags (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			is_active BOOLEAN NOT NULL DEFAULT FALSE
		);
	`)
}

func setupTagToItemSchema(db *sqlx.DB) error {
	return setupTable(db, "tag_to_item", `
		CREATE TABLE tag_to_item (
			tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
			item_id UUID REFERENCES items(id) ON DELETE CASCADE,

			PRIMARY KEY (item_id, tag_id)
		);
	`)
}

func setupTable(db *sqlx.DB, name string, schema string) error {
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create %s schema: %w", name, err)
	}

	return nil
}
