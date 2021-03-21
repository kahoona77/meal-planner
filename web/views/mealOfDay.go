package views

import (
	"database/sql"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"net/http"
	"strconv"
)

func MealOfDay(ctx *core.WebContext) error {
	repo := planner.NewRepository(ctx)
	mealsRepo := meals.NewRepository(ctx)

	id := ctx.Param("id")
	mealOfDay, err := repo.GetMealOfDay(id)
	if err != nil {
		mealOfDay = planner.NewMealFromId(id)
	} else {
		meal, merr := mealsRepo.GetMeal(mealOfDay.MealId.Int64)
		if merr != nil {
			mealOfDay.Meal = &meals.Meal{}
		} else {
			mealOfDay.Meal = meal
		}
	}

	return ctx.Render(http.StatusOK, "meal-of-day.html", mealOfDay)
}

func SelectMealOfDayView(ctx *core.WebContext) error {
	repo := planner.NewRepository(ctx)
	mealsRepo := meals.NewRepository(ctx)

	id := ctx.Param("id")
	mealOfDay, err := repo.GetMealOfDay(id)
	if err != nil {
		mealOfDay = planner.NewMealFromId(id)
	}

	allMeals, err := mealsRepo.GetMeals()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"mealOfDay": mealOfDay,
		"meals":     allMeals,
	}

	return ctx.Render(http.StatusOK, "meal-of-day-select.html", data)
}

func SelectMealOfDay(ctx *core.WebContext) error {
	repo := planner.NewRepository(ctx)

	id := ctx.Param("id")
	mealOfDay, err := repo.GetMealOfDay(id)
	createNew := false
	if err != nil {
		mealOfDay = planner.NewMealFromId(id)
		createNew = true
	}

	selected := ctx.FormValue("selected")
	mealId, err := strconv.Atoi(selected)
	if err != nil {
		mealOfDay.MealId = sql.NullInt64{}
	} else {
		mealOfDay.MealId = sql.NullInt64{Int64: int64(mealId), Valid: true}
	}

	if createNew {
		if err := repo.CreateMealOfDay(mealOfDay); err != nil {
			return err
		}
	} else {
		if err := repo.UpdateMealOfDay(mealOfDay); err != nil {
			return err
		}
	}

	return ctx.Redirect(http.StatusFound, "/")
}
