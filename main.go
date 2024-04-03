package main

import (
	"embed"
	"fmt"
	"meal-planner/core"
	"meal-planner/files"
	"meal-planner/web"
	"meal-planner/web/views"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	app := core.InitApp(web.CreateRenderer, embedMigrations)
	defer app.Ctx.Close()

	//files.InitDb(app.Ctx)
	//meals.InitDb(app.Ctx)
	//planner.InitDb(app.Ctx)

	root := app.Group(app.Ctx.Config().BasePath)

	root.Static("/assets", "./web/assets/dist")

	root.GET("/files/:id", files.GetFile)

	root.GET("/", views.Index)
	root.GET("/offset/:offset", views.Offset)

	root.GET("/wizard/:offset", views.Wizard)

	root.GET("/planner/:id", views.MealOfDay)
	root.GET("/planner/:id/select", views.SelectMealOfDayView)
	root.POST("/planner/:id/select", views.SelectMealOfDay)

	root.GET("/meals", views.Meals)
	root.POST("/meals", views.MealSave)
	root.GET("/meals/new", views.MealEdit)
	root.GET("/meals/:id", views.MealEdit)
	root.POST("/meals/:id", views.MealSave)
	root.POST("/meals/:id/delete", views.MealDelete)

	root.GET("/tags", views.Tags)
	root.POST("/tags", views.TagSave)
	root.POST("/tags/:id", views.TagSave)

	// Listen and server on 0.0.0.0:8080
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", app.Ctx.Config().Port)))
}
