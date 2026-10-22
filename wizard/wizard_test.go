package wizard

import (
	"meal-planner/internal/testdb"
	"meal-planner/meals"
	"meal-planner/planner"
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

func TestReroll(t *testing.T) {
	db := testdb.New(t)
	db.MustExec(`INSERT INTO tags (id, name) VALUES (1, 'Suppe'), (2, 'Nudeln')`)
	db.MustExec(`INSERT INTO meals (id, name) VALUES (1, 'Tomatensuppe'), (2, 'Kürbissuppe'), (3, 'Spaghetti'), (4, 'Lasagne')`)
	db.MustExec(`INSERT INTO meal_tags (meal_id, tag_id) VALUES (1, 1), (2, 1), (3, 2), (4, 2)`)
	monday := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	wizard := NewWizard(testdb.Context{DB: db})

	planned := func(mealIds ...int64) []*planner.MealOfTheDay {
		mods := make([]*planner.MealOfTheDay, len(mealIds))
		for i, id := range mealIds {
			mods[i] = planner.NewMeal(monday.AddDate(0, 0, i))
			mods[i].SetMeal(&meals.Meal{Id: id})
		}
		return mods
	}
	soup := &Day{Weekday: 1, Date: monday, Tags: []*WeekdayTag{{TagId: 1}}}

	t.Run("picks another meal with the tags of the day", func(t *testing.T) {
		mod, err := wizard.Reroll(soup, planned(1))
		require.NoError(t, err)
		assert.Equal(t, "Kürbissuppe", mod.Meal.Name)

		var mealId int64
		require.NoError(t, db.Get(&mealId, `SELECT meal_id FROM meals_of_day WHERE id = '2026_10_05'`))
		assert.Equal(t, int64(2), mealId)
	})

	t.Run("ignores the tags if all meals with them are planned", func(t *testing.T) {
		mod, err := wizard.Reroll(soup, planned(1, 2, 3))
		require.NoError(t, err)
		assert.Equal(t, "Lasagne", mod.Meal.Name)
	})

	t.Run("fails if every meal is planned", func(t *testing.T) {
		_, err := wizard.Reroll(soup, planned(1, 2, 3, 4))
		assert.Error(t, err)
	})
}
