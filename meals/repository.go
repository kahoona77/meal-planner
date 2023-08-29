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
	result, err := r.db.NamedExec(mealsInsert, meal)
	if err != nil {
		return err
	}

	// set inserted Id
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	meal.Id = id

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

func (r *Repository) GetTags() ([]*Tag, error) {
	var tags []*Tag
	err := r.db.Select(&tags, "SELECT * FROM tags")
	return tags, err
}

func (r *Repository) GetTag(id int64) (*Tag, error) {
	var tag Tag
	err := r.db.Get(&tag, "SELECT * FROM tags WHERE id=$1", id)
	return &tag, err
}

const tagsInsert = `INSERT INTO tags (name, color) VALUES (:name, :color)`

func (r *Repository) CreateTag(tag *Tag) error {
	_, err := r.db.NamedExec(tagsInsert, tag)
	return err
}

const tagsUpdate = `UPDATE tags SET name=:name, color=:color WHERE id = :id`

func (r *Repository) UpdateTag(tag *Tag) error {
	_, err := r.db.NamedExec(tagsUpdate, tag)
	return err
}

func (r *Repository) GetMealTags(mealId int64) ([]*MealTag, error) {
	var mealTags []*MealTag
	err := r.db.Select(&mealTags, "SELECT * FROM meal_tags left join tags t on t.id = meal_tags.tag_id where meal_id=$1", mealId)
	return mealTags, err
}

const mealTagsInsert = `INSERT INTO meal_tags (meal_id, tag_id) VALUES (:meal_id, :tag_id)`

func (r *Repository) SetTagsForMeal(meal *Meal, tags []*Tag) error {
	// delete old tags
	_, err := r.db.Exec("DELETE FROM meal_tags WHERE meal_id=$1", meal.Id)
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		mealTags := make([]*MealTag, len(tags))
		for i, tag := range tags {
			mealTags[i] = &MealTag{
				MealId: meal.Id,
				TagId:  tag.Id,
			}
		}

		// insert new tags
		_, err = r.db.NamedExec(mealTagsInsert, mealTags)
		return err
	}

	return nil
}
