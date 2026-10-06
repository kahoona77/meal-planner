package views

import (
	"log/slog"
	"meal-planner/core"
	"meal-planner/meals"
	"net/http"
)

func Tags(ctx *core.WebContext) error {
	repo := meals.NewRepository(ctx)
	tags, err := repo.GetTags()
	if err != nil {
		slog.Error("error loading tags", "err", err)
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
			return err
		}
	} else {
		if err := repo.UpdateTag(tag); err != nil {
			return err
		}
	}

	return ctx.Redirect(http.StatusFound, "/tags")
}
