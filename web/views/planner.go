package views

import (
	"github.com/labstack/echo/v4"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"net/http"
	"time"
)

func Index(c echo.Context) error {
	ctx := c.(*core.WebContext)
	return renderWeek(ctx, 0)
}

func Offset(c echo.Context) error {
	ctx := c.(*core.WebContext)
	offset := ctx.ParamAsInt("offset")
	return renderWeek(ctx, offset)
}

func renderWeek(ctx *core.WebContext, offset int) error {
	week, err := getWeek(ctx, offset)
	if err != nil {
		return err
	}

	return ctx.Render(http.StatusOK, "index.html", &week)
}

func getWeek(ctx core.Context, offset int) (planner.Week, error) {
	now := time.Now().AddDate(0, 0, offset*7)
	week := planner.Week{
		Start:  planner.GetStartWeek(now),
		End:    planner.GetEndWeek(now),
		Offset: offset,
	}

	week.Meals = make([]*planner.MealOfTheDay, 7)

	repo := planner.NewRepository(ctx)
	mealRepo := meals.NewRepository(ctx)
	mealsOfTheDay, err := repo.GetMealsOfDay(week.Start, week.End)
	if err != nil {
		return week, err
	}

	date := week.Start
	for i := 0; i < 7; i++ {

		meal, err := findMeal(mealRepo, date, mealsOfTheDay)
		if err != nil {
			return week, err
		}

		week.Meals[i] = meal
		date = date.AddDate(0, 0, 1)
	}
	return week, nil
}

func findMeal(mealRepo *meals.Repository, date time.Time, mealsOfTheDay []*planner.MealOfTheDay) (*planner.MealOfTheDay, error) {
	var result *planner.MealOfTheDay = nil
	for _, meal := range mealsOfTheDay {
		dateId := date.Format(planner.MealOfDayDateFormat)
		if dateId == meal.Id {
			result = meal
		}
	}

	if result == nil {
		result = &planner.MealOfTheDay{
			Date: date,
			Meal: &meals.Meal{},
		}
	} else {
		meal, err := mealRepo.GetMeal(result.MealId)
		if err != nil {
			return nil, err
		}
		result.Meal = meal
	}

	result.UpdateId()
	return result, nil
}
