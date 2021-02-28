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
	//ctx := c.(*core.WebContext).Ctx

	week := getWeek(0)

	return c.Render(http.StatusOK, "index.html", &week)
}

func Offset(c echo.Context) error {
	ctx := c.(*core.WebContext)
	offset := ctx.ParamAsInt("offset")

	week := getWeek(offset)

	return c.Render(http.StatusOK, "index.html", &week)
}

func getWeek(offset int) planner.Week {
	now := time.Now().AddDate(0, 0, offset*7)
	week := planner.Week{
		Start:  planner.GetStartWeek(now),
		End:    planner.GetEndWeek(now),
		Offset: offset,
	}

	week.Meals = make([]planner.MealOfTheDay, 7)

	date := week.Start
	for i := 0; i < 7; i++ {
		week.Meals[i] = planner.MealOfTheDay{Date: date, Meal: meals.Meal{Id: 0, Name: "Fischstäbchen mit Kartoffeln"}}
		date = date.AddDate(0, 0, 1)
	}
	return week
}
