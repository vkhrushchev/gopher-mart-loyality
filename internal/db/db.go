package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ExecuteMigrations(databaseDSN string) error {
	m, err := migrate.New("file://../../internal/db/migrations", databaseDSN)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
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
