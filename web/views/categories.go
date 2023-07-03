package views

import (
	"github.com/sirupsen/logrus"
	"meal-planner/core"
	"meal-planner/meals"
	"net/http"
)

func Categories(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)
	categories, err := repo.GetCategories()
	if err != nil {
		logrus.Errorf("error loading categories: %v", err)
	}

	data := core.TemplateData{
		"categories": categories,
	}

	return ctx.RenderTemplate(http.StatusOK, "category-list.html", data)
}

func CategorySave(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)

	category := &meals.Category{}
	isNew := ctx.Param("id") == ""
	if !isNew {
		var err error
		category, err = repo.GetCategory(int64(ctx.ParamAsInt("id")))
		if err != nil {
			return err
		}
	}

	//update
	category.Name = ctx.FormValue("name")

	if isNew {
		if err := repo.CreateCategory(category); err != nil {
			logrus.Errorf("error creating meal: %v", err)
			return err
		}
	} else {
		if err := repo.UpdateCategory(category); err != nil {
			logrus.Errorf("error updating meal: %v", err)
			return err
		}
	}

	return ctx.Redirect(http.StatusFound, "/categories")
}
