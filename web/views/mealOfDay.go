package views

import (
	"database/sql"
	"log/slog"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"meal-planner/wizard"
	"net/http"
	"strconv"
	"strings"
	"time"
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

// RerollMealOfDay replaces the meal of a day by another random one, or picks one for a day
// without a meal.
func RerollMealOfDay(ctx *core.WebContext) error {
	repo := planner.NewRepository(ctx)

	id := ctx.Param("id")
	if _, err := time.Parse(planner.MealOfDayDateFormat, id); err != nil {
		// not a day: 404
		return sql.ErrNoRows
	}
	mealOfDay := planner.NewMealFromId(id)

	mealsOfWeek, err := repo.GetMealsOfDay(planner.GetStartWeek(mealOfDay.Date), planner.GetEndWeek(mealOfDay.Date))
	if err != nil {
		return err
	}
	for _, mod := range mealsOfWeek {
		mod.Meal = &meals.Meal{Id: mod.MealId.Int64}
	}

	weekdayTags, err := wizard.NewRepository(ctx).GetWeekdayTags()
	if err != nil {
		return err
	}
	weekday := int(mealOfDay.Date.Weekday())
	day := &wizard.Day{Weekday: weekday, Date: mealOfDay.Date, Tags: getWeekdayTags(weekday, weekdayTags)}

	if _, err := wizard.NewWizard(ctx).Reroll(day, mealsOfWeek); err != nil {
		// e.g. every meal is planned this week already: the day keeps its meal
		slog.Warn("could not reroll the meal of the day", "id", mealOfDay.Id, "err", err)
	}

	return ctx.Redirect(http.StatusFound, "/"+strings.TrimPrefix(weekURL(mealOfDay.Date), "./"))
}
