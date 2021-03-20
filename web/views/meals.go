package views

import (
	"github.com/sirupsen/logrus"
	"meal-planner/core"
	"meal-planner/meals"
	"net/http"
)

func Meals(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)
	mealsList, err := repo.GetMeals()
	if err != nil {
		logrus.Errorf("error loading meals: %v", err)
	}

	data := map[string]interface{}{
		"meals": mealsList,
	}

	return ctx.Render(http.StatusOK, "meals-list.html", data)
}

func MealEdit(ctx *core.WebContext) error {
	meal := &meals.Meal{}
	id := ctx.Param("id")

	if id != "" {
		var err error
		meal, err = meals.NewRepository(ctx.Ctx).GetMeal(ctx.ParamAsInt("id"))
		if err != nil {
			return err
		}
	}

	return ctx.Render(http.StatusOK, "meals-edit.html", meal)
}

func MealSave(ctx *core.WebContext) error {
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
	meal.Name = ctx.FormValue("name")
	meal.Description = ctx.FormValue("description")

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

func MealDelete(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)

	id := ctx.ParamAsInt("id")
	if err := repo.DeleteMeal(id); err != nil {
		return err
	}

	return ctx.Redirect(http.StatusFound, "/meals")
}
