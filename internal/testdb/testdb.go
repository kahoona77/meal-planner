// Package testdb provides an in-memory database with all migrations applied, for tests.
package testdb

import (
	"meal-planner/core"
	"meal-planner/migrations"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// New returns an empty in-memory database with all migrations applied.
// It is closed when the test ends.
func New(t *testing.T) *sqlx.DB {
	t.Helper()
	db := sqlx.MustConnect("sqlite3", ":memory:")
	// every connection would get its own in-memory database
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })

	if err := migrations.Up(db.DB); err != nil {
		t.Fatal(err)
	}
	return db
}

// Context is a core.Context for repositories in tests.
type Context struct {
	DB *sqlx.DB
}

func (c Context) Db() *sqlx.DB            { return c.DB }
func (c Context) Config() *core.AppConfig { return &core.AppConfig{} }
