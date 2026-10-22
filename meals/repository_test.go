package meals

import (
	"meal-planner/internal/testdb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteMeal_removesTags(t *testing.T) {
	db := testdb.New(t)
	db.MustExec(`INSERT INTO tags (id, name) VALUES (1, 'Suppe')`)
	db.MustExec(`INSERT INTO meals (id, name) VALUES (1, 'Tomatensuppe'), (2, 'Kürbissuppe')`)
	db.MustExec(`INSERT INTO meal_tags (meal_id, tag_id) VALUES (1, 1), (2, 1)`)

	require.NoError(t, NewRepository(testdb.Context{DB: db}).DeleteMeal(1))

	var mealIds []int64
	require.NoError(t, db.Select(&mealIds, `SELECT meal_id FROM meal_tags`))
	assert.Equal(t, []int64{2}, mealIds)
}
