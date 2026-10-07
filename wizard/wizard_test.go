package wizard

import (
	"meal-planner/core"
	"os"
	"testing"
	"time"

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

func TestGenerate_skipsDaysWithoutMeal(t *testing.T) {
	db := newTestDb(t)
	db.MustExec(`INSERT INTO tags (id, name) VALUES (1, 'Suppe'), (2, 'Nudeln')`)
	db.MustExec(`INSERT INTO meals (id, name) VALUES (1, 'Spaghetti')`)
	// tag 1 only belongs to a meal that does not exist anymore
	db.MustExec(`INSERT INTO meal_tags (meal_id, tag_id) VALUES (0, 1), (1, 2)`)

	monday := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	week := Week{
		Start: monday,
		End:   monday.AddDate(0, 0, 7),
		Days: []*Day{
			{Weekday: 1, Date: monday, Tags: []*WeekdayTag{{TagId: 1}}},
			{Weekday: 2, Date: monday.AddDate(0, 0, 1), Tags: []*WeekdayTag{{TagId: 2}}},
		},
	}

	plannerWeek, err := NewWizard(testContext{db: db}).Generate(week)

	require.NoError(t, err)
	require.Len(t, plannerWeek.Meals, 1)
	assert.Equal(t, "Spaghetti", plannerWeek.Meals[0].Meal.Name)
}
