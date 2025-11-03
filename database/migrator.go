package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(databaseUrl string) {
	m, err := migrate.New(
		"file://database/postgres",
		databaseUrl,
	)
	if err != nil {
		log.Fatalf("migrate init error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up error: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Printf("no changes to migrate")
	} else {
		log.Println("migrations applied successfully")
	}
}
