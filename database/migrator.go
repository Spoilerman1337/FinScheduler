package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(dsn string) {
	m, err := migrate.New(
		"file://database/postgres",
		dsn,
	)
	if err != nil {
		log.Fatalf("migrate init error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up error: %v", err)
	}

	log.Println("migrations applied successfully")
}
