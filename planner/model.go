package planner

import (
	"meal-planner/core"
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
	Id     string      `db:"id"`
	Date   time.Time   `db:"date"`
	MealId int         `db:"meal_id"`
	Meal   *meals.Meal `db:"-"`
}

func (m *MealOfTheDay) UpdateId() {
	m.Id = m.Date.Format(MealOfDayDateFormat)
}

func InitDb(ctx *core.Ctx) {
	ctx.Db().MustExec(plannerSchema)
}

var plannerSchema = `
CREATE TABLE IF NOT EXISTS meals_of_day (
	id          VARCHAR(10) PRIMARY KEY,
	date        TIMESTAMP NOT NULL,
	meal_id     INTEGER NOT NULL
);`
