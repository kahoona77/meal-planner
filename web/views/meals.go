package views

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"meal-planner/core"
	"meal-planner/meals"
	"net/http"
)

func Meals(c echo.Context) error {
	ctx := c.(*core.WebContext)

	repo := meals.NewRepository(ctx)
	mealsList, err := repo.GetMeals()
	if err != nil {
		logrus.Errorf("error loading meals: %v", err)
	}

	data := map[string]interface{}{
		"meals": mealsList,
	}

	return c.Render(http.StatusOK, "meals-list.html", data)
}

func MealEdit(c echo.Context) error {
	ctx := c.(*core.WebContext)

	meal := &meals.Meal{}
	id := ctx.Param("id")

	if id != "" {
		var err error
		meal, err = meals.NewRepository(ctx.Ctx).GetMeal(ctx.ParamAsInt("id"))
		if err != nil {
			return err
		}
	}

	return c.Render(http.StatusOK, "meals-edit.html", meal)
}

func MealSave(c echo.Context) error {
	ctx := c.(*core.WebContext)
	repo := meals.NewRepository(ctx)

	meal := &meals.Meal{}
	isNew := ctx.Param("id") == ""
	if !isNew {
		var err error
		meal, err = meals.NewRepository(ctx.Ctx).GetMeal(ctx.ParamAsInt("id"))
		if err != nil {
			return err
		}
	}

	//update
	meal.Name = c.FormValue("name")
	meal.Description = c.FormValue("description")

	if isNew {
		if err := repo.CreateMeal(meal); err != nil {
			logrus.Errorf("error creating meal: %v", err)
			return err
		}
	} else {
		if err := repo.UpdateMeal(meal); err != nil {
			logrus.Errorf("error updating meal: %v", err)
			return err
		}
	}

	return ctx.Redirect(http.StatusFound, "/meals")
}

func MealDelete(c echo.Context) error {
	ctx := c.(*core.WebContext)
	repo := meals.NewRepository(ctx)

	id := ctx.ParamAsInt("id")
	if err := repo.DeleteMeal(id); err != nil {
		return err
	}

	return ctx.Redirect(http.StatusFound, "/meals")
}
