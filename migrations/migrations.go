// Package migrations holds the database schema as goose SQL migrations and applies them.
package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var files embed.FS

// Up applies all pending migrations to db.
func Up(db *sql.DB) error {
	goose.SetBaseFS(files)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("error selecting migration dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
