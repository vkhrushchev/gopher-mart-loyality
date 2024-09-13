package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

//go:embed migrations/*.sql
var migrationsDir embed.FS

func ExecuteMigrations(databaseDSN string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseDSN)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func NewSqlxDB(databaseDSN string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("db: error when open database: %w", err)
	}

	return db, nil
}
