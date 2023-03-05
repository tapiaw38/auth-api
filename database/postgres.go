package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new postgres repository
func NewPostresRepository(url string) (*PostgresRepository, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

// GetRelativePathToMigrationsDirectory gets the relative path to migrations directory
func GetRelativePathToMigrationsDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absMigrationsDirPath := filepath.Join(cwd, "database", "migrations")

	relMigrationsDirPath, err := filepath.Rel(cwd, absMigrationsDirPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("file://%s", relMigrationsDirPath), nil
}

// Makemigration makes the database migrations
func (postgres *PostgresRepository) Makemigration(databaseUrl string) error {
	migrationPath, err := GetRelativePathToMigrationsDirectory()
	if err != nil {
		return err
	}

	m, err := migrate.New(migrationPath, databaseUrl)
	if err != nil {
		return err
	}

	version, _, _ := m.Version()
	log.Printf("migrations: current version is %v", version)

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("migrations:", err)
			return nil
		}
		return err
	}

	log.Println("migrations: database migrated")

	return nil
}

// CheckConnection checks the database connection
func (postgres *PostgresRepository) CheckConnection() bool {
	db, err := postgres.db.Conn(context.Background())
	if err != nil {
		panic(err)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		panic(err)
	}

	return err == nil
}
