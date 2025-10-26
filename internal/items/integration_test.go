//go:build integration
// +build integration

package items

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, terminate, db, err := setupPostgresContainer(ctx)
	if err != nil {
		fmt.Printf("failed to start postgres: %v", err)
	}

	testDB = db
	//terminateContainer = terminate

	schema := `
        CREATE TABLE items (
		   id UUID PRIMARY KEY,
		   name TEXT NOT NULL UNIQUE,
		   price NUMERIC(16, 2) NOT NULL DEFAULT 0 CHECK (price >= 0),
		   description TEXT NULL,
		   is_active BOOLEAN NOT NULL DEFAULT FALSE,
		   created_at TIMESTAMP NOT NULL DEFAULT now(),
		   updated_at TIMESTAMP NULL
		);
    `

	_, err = testDB.Exec(schema)
	if err != nil {
		fmt.Printf("failed to create schema: %v", err)
	}

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
