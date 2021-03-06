package meals

import (
	"github.com/jmoiron/sqlx"
	"meal-planner/core"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(c core.Context) *Repository {
	return &Repository{db: c.Db()}
}

const mealsInsert = `INSERT INTO meals (name, description, image_file_id) VALUES (:name, :description, :image_file_id)`

func (r *Repository) CreateMeal(meal *Meal) error {
	_, err := r.db.NamedExec(mealsInsert, meal)
	return err
}

const mealsUpdate = `UPDATE meals SET name=:name, description=:description, image_file_id=:image_file_id WHERE id = :id`

func (r *Repository) UpdateMeal(meal *Meal) error {
	_, err := r.db.NamedExec(mealsUpdate, meal)
	return err
}

func (r *Repository) GetMeal(id int64) (*Meal, error) {
	var meal Meal
	err := r.db.Get(&meal, "SELECT * FROM meals WHERE id=$1", id)
	return &meal, err
}

func (r *Repository) GetMeals() ([]*Meal, error) {
	var meals []*Meal
	err := r.db.Select(&meals, "SELECT * FROM meals")
	return meals, err
}

func (r *Repository) DeleteMeal(id int64) error {
	_, err := r.db.Exec("DELETE FROM meals WHERE id=$1", id)
	return err
}
