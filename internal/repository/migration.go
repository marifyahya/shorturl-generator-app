package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName, driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %v", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

func RollbackMigration(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName, driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %v", err)
	}

	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not rollback migration: %v", err)
	}

	log.Println("⚠️ Database rollback (1 step) completed successfully")
	return nil
}

