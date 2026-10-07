package meals

import (
	"meal-planner/core"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testContext struct {
	db *sqlx.DB
}

func (c testContext) Db() *sqlx.DB            { return c.db }
func (c testContext) Config() *core.AppConfig { return &core.AppConfig{} }

// newTestDb returns an in-memory database with all migrations applied.
func newTestDb(t *testing.T) *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	// every connection would get its own in-memory database
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })

	goose.SetBaseFS(os.DirFS(".."))
	require.NoError(t, goose.SetDialect("sqlite"))
	require.NoError(t, goose.Up(db.DB, "migrations"))
	return db
}

func TestDeleteMeal_removesTags(t *testing.T) {
	db := newTestDb(t)
	db.MustExec(`INSERT INTO tags (id, name) VALUES (1, 'Suppe')`)
	db.MustExec(`INSERT INTO meals (id, name) VALUES (1, 'Tomatensuppe'), (2, 'Kürbissuppe')`)
	db.MustExec(`INSERT INTO meal_tags (meal_id, tag_id) VALUES (1, 1), (2, 1)`)

	require.NoError(t, NewRepository(testContext{db: db}).DeleteMeal(1))

	var mealIds []int64
	require.NoError(t, db.Select(&mealIds, `SELECT meal_id FROM meal_tags`))
	assert.Equal(t, []int64{2}, mealIds)
}
