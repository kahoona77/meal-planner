package main

import (
	"fmt"
	"log/slog"
	"meal-planner/core"
	"meal-planner/files"
	"meal-planner/web"
	"meal-planner/web/views"
	"os"
)

func main() {
	app := core.InitApp()
	defer app.Ctx.Close()

	setupRoutes(app)

	if err := app.Start(fmt.Sprintf(":%s", app.Ctx.Config().Port)); err != nil {
		slog.Error("http server stopped", "err", err)
		os.Exit(1)
	}
}

// setupRoutes registers the middlewares and all routes of the app.
func setupRoutes(app *core.App) {
	basePath := app.Ctx.Config().BasePath
	static := web.StaticFS(app.Ctx.Config().StaticDir)
	app.Use(web.PageMiddleware(basePath, web.NewAssets(static, basePath)))

	root := app.Group(basePath)

	root.Static("/assets", static)

	root.GET("/files/{id}", files.GetFile)

	root.GET("/", views.Index)
	root.GET("/offset/{offset}", views.Offset)

	root.GET("/wizard/{offset}", views.Wizard)
	root.POST("/wizard/{offset}", views.Generate)

	root.GET("/planner/{id}", views.MealOfDay)
	root.GET("/planner/{id}/select", views.SelectMealOfDayView)
	root.POST("/planner/{id}/select", views.SelectMealOfDay)
	root.POST("/planner/{id}/reroll", views.RerollMealOfDay)

	root.GET("/meals", views.Meals)
	root.POST("/meals", views.MealSave)
	root.GET("/meals/new", views.MealEdit)
	root.GET("/meals/{id}", views.MealEdit)
	root.POST("/meals/{id}", views.MealSave)
	root.POST("/meals/{id}/delete", views.MealDelete)

	root.GET("/tags", views.Tags)
	root.POST("/tags", views.TagSave)
	root.POST("/tags/{id}", views.TagSave)
}
