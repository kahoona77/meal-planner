package views

import (
	"database/sql"
	"meal-planner/planner"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekURL(t *testing.T) {
	now := time.Now()

	assert.Equal(t, "./", weekURL(now))
	assert.Equal(t, "offset/-1", weekURL(now.AddDate(0, 0, -7)))
	assert.Equal(t, "offset/3", weekURL(now.AddDate(0, 0, 21)))
	// any day of a week leads to it
	for i := range 7 {
		assert.Equal(t, "offset/1", weekURL(planner.GetStartWeek(now).AddDate(0, 0, 7+i)))
	}
}

func TestSelectBackURL(t *testing.T) {
	mod := planner.NewMeal(time.Now().AddDate(0, 0, -7))
	assert.Equal(t, "offset/-1", selectBackURL(mod))

	mod.MealId = sql.NullInt64{Int64: 1, Valid: true}
	assert.Equal(t, "planner/"+mod.Id, selectBackURL(mod))
}
