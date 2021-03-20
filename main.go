package main

import (
	"fmt"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"meal-planner/web/views"
)

func main() {
	app := core.InitApp()
	defer app.Ctx.Close()

	meals.InitDb(app.Ctx)
	planner.InitDb(app.Ctx)

	root := app.Group(app.Ctx.Config().BasePath)

	root.Static("/assets", "./web/assets")

	root.GET("/", views.Index)
	root.GET("/offset/:offset", views.Offset)

	root.GET("/planner/:id", views.MealOfDay)
	root.GET("/planner/:id/select", views.SelectMealOfDayView)
	root.POST("/planner/:id/select", views.SelectMealOfDay)

	root.GET("/meals", views.Meals)
	root.POST("/meals", views.MealSave)
	root.GET("/meals/new", views.MealEdit)
	root.GET("/meals/:id", views.MealEdit)
	root.POST("/meals/:id", views.MealSave)
	root.POST("/meals/:id/delete", views.MealDelete)

	// Listen and server on 0.0.0.0:8080
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", app.Ctx.Config().Port)))
}
