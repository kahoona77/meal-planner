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

	var tags []*meals.MealTag
	if mealOfDay.Meal.Id > 0 {
		tags, err = mealsRepo.GetMealTags(mealOfDay.Meal.Id)
		if err != nil {
			return err
		}
	}

	return ctx.Render(http.StatusOK, mealOfDayPage(mealOfDay, tags))
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

	return ctx.Render(http.StatusOK, selectMealOfDayPage(mealOfDay, allMeals))
}

func SelectMealOfDay(ctx *core.WebContext) error {
	repo := planner.NewRepository(ctx)

	id := ctx.Param("id")
	mealOfDay, err := repo.GetMealOfDay(id)
	if err != nil {
		mealOfDay = planner.NewMealFromId(id)
	}

	selected := ctx.FormValue("selected")
	mealId, err := strconv.Atoi(selected)
	if err != nil {
		mealOfDay.MealId = sql.NullInt64{}
	} else {
		mealOfDay.MealId = sql.NullInt64{Int64: int64(mealId), Valid: true}
	}

	if err := repo.UpsertMealOfDay(mealOfDay); err != nil {
		return err
	}

	return ctx.Redirect(http.StatusFound, "/")
}
