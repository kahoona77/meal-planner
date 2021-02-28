package planner

import (
	"meal-planner/meals"
	"time"
)

type Week struct {
	Start  time.Time
	End    time.Time
	Offset int
	Meals  []MealOfTheDay
}

type MealOfTheDay struct {
	Date time.Time
	Meal meals.Meal
}
