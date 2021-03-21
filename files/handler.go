package files

import (
	"meal-planner/core"
	"net/http"
)

func GetFile(ctx *core.WebContext) error {
	repo := NewRepository(ctx)

	id := ctx.ParamAsInt("id")
	file, err := repo.GetFile(int64(id))
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusOK, file.ContentType, file.Data)
}
