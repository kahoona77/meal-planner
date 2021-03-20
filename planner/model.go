package planner

import (
	"github.com/sirupsen/logrus"
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

	repo := NewRepository(ctx)

	date, _ := time.Parse(MealOfDayDateFormat, "2021_03_20")

	if err := repo.CreateMealOfDay(&MealOfTheDay{
		Id:     "2021_03_20",
		Date:   date,
		MealId: 1,
	}); err != nil {
		logrus.Errorf("err: %v", err)
	}

}

var plannerSchema = `
CREATE TABLE IF NOT EXISTS meals_of_day (
	id          VARCHAR(10) PRIMARY KEY,
	date        TIMESTAMP NOT NULL,
	meal_id     INTEGER NOT NULL
);`
