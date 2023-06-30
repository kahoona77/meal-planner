package core

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"io/fs"
)

func runMigrations(db *sql.DB, migrationsFs fs.FS) error {
	goose.SetBaseFS(migrationsFs)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("error selecting migration dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
