package views

import (
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"net/http"
)

func MealOfDay(ctx *core.WebContext) error {
	repo := planner.NewRepository(ctx)
	mealsRepo := meals.NewRepository(ctx)

	id := ctx.Param("id")
	mealOfDay, err := repo.GetMealOfDay(id)
	if err != nil {
		mealOfDay = &planner.MealOfTheDay{
			Meal: &meals.Meal{},
		}
	} else {
		meal, merr := mealsRepo.GetMeal(mealOfDay.MealId)
		if merr != nil {
			return merr
		}
		mealOfDay.Meal = meal
	}

	return ctx.Render(http.StatusOK, "meal-of-day.html", mealOfDay)
}
