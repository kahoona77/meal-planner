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
