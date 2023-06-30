package planner

import (
	"database/sql"
	"meal-planner/meals"
	"time"
)

const MealOfDayDateFormat = "2006_01_02"

type Week struct {
	Start  time.Time
	End    time.Time
	Offset int
	Meals  []*MealOfTheDay
}

func NewMealFromId(id string) *MealOfTheDay {
	date, _ := time.Parse(MealOfDayDateFormat, id)
	m := &MealOfTheDay{Date: date, Meal: &meals.Meal{}}
	m.UpdateId()
	return m
}

func NewMeal(date time.Time) *MealOfTheDay {
	m := &MealOfTheDay{Date: date, Meal: &meals.Meal{}}
	m.UpdateId()
	return m
}

type MealOfTheDay struct {
	Id     string        `db:"id"`
	Date   time.Time     `db:"date"`
	MealId sql.NullInt64 `db:"meal_id"`
	Meal   *meals.Meal   `db:"-"`
}

func (m *MealOfTheDay) UpdateId() {
	m.Id = m.Date.Format(MealOfDayDateFormat)
}
