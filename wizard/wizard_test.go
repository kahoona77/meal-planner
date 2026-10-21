package wizard

import (
	"meal-planner/internal/testdb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_skipsDaysWithoutMeal(t *testing.T) {
	db := testdb.New(t)
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

	plannerWeek, err := NewWizard(testdb.Context{DB: db}).Generate(week)

	require.NoError(t, err)
	require.Len(t, plannerWeek.Meals, 1)
	assert.Equal(t, "Spaghetti", plannerWeek.Meals[0].Meal.Name)
}
