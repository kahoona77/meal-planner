package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"meal-planner/core"
	"meal-planner/meals"
	"meal-planner/planner"
	"meal-planner/web"
	"meal-planner/web/views"
)

func main() {
	ctx := core.InitApp()
	defer ctx.Close()

	meals.InitDb(ctx)
	planner.InitDb(ctx)

	e := echo.New()
	e.Renderer = web.NewTemplate(ctx)
	e.Debug = true
	//e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(core.CreateCtx(ctx))

	root := e.Group(ctx.Config().BasePath)

	root.Static("/assets", "./web/assets")

	root.GET("/", views.Index)
	root.GET("/offset/:offset", views.Offset)

	root.GET("/planner/:id", views.MealOfDay)

	root.GET("/meals", views.Meals)
	root.POST("/meals", views.MealSave)
	root.GET("/meals/new", views.MealEdit)
	root.GET("/meals/:id", views.MealEdit)
	root.POST("/meals/:id", views.MealSave)
	root.POST("/meals/:id/delete", views.MealDelete)

	// Listen and server on 0.0.0.0:8080
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", ctx.Config().Port)))
}
