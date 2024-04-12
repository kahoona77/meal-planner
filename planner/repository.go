package planner

import (
	"github.com/jmoiron/sqlx"
	"meal-planner/core"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(c core.Context) *Repository {
	return &Repository{db: c.Db()}
}

const mealsInsert = `INSERT INTO meals_of_day (id, date, meal_id) VALUES (:id, :date, :meal_id)`

const mealsUpsert = `INSERT INTO meals_of_day (id, date, meal_id) VALUES (:id, :date, :meal_id) ON CONFLICT (id) DO UPDATE SET date=:date, meal_id=:meal_id WHERE id = :id`

func (r *Repository) UpsertMealOfDay(meal *MealOfTheDay) error {
	meal.UpdateId()
	_, err := r.db.NamedExec(mealsUpsert, meal)
	return err
}

func (r *Repository) GetMealOfDay(mealId string) (*MealOfTheDay, error) {
	var meal MealOfTheDay
	err := r.db.Get(&meal, "SELECT * FROM meals_of_day WHERE id=$1", mealId)
	return &meal, err
}

func (r *Repository) GetMealsOfDay(start time.Time, end time.Time) ([]*MealOfTheDay, error) {
	var meals []*MealOfTheDay
	err := r.db.Select(&meals, "SELECT * FROM meals_of_day where date >= $1 AND date < $2  ORDER BY date", start.UTC(), end.UTC())
	return meals, err
}

func (r *Repository) DeleteMealOfDay(mealId string) error {
	_, err := r.db.Exec("DELETE FROM meals_of_day WHERE id=$1", mealId)
	return err
}
