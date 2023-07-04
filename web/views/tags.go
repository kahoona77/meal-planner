package views

import (
	"github.com/sirupsen/logrus"
	"meal-planner/core"
	"meal-planner/meals"
	"net/http"
)

func Tags(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)
	tags, err := repo.GetTags()
	if err != nil {
		logrus.Errorf("error loading tags: %v", err)
	}

	data := core.TemplateData{
		"tags": tags,
	}

	return ctx.RenderTemplate(http.StatusOK, "tag-list.html", data)
}

func TagSave(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)

	tag := &meals.Tag{}
	isNew := ctx.Param("id") == ""
	if !isNew {
		var err error
		tag, err = repo.GetTag(int64(ctx.ParamAsInt("id")))
		if err != nil {
			return err
		}
	}

	//update
	tag.Name = ctx.FormValue("name")
	tag.Color = ctx.FormValue("color")

	if isNew {
		if err := repo.CreateTag(tag); err != nil {
			logrus.Errorf("error creating meal: %v", err)
			return err
		}
	} else {
		if err := repo.UpdateTag(tag); err != nil {
			logrus.Errorf("error updating meal: %v", err)
			return err
		}
	}

	return ctx.Redirect(http.StatusFound, "/tags")
}
