package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"meal-planner/core"
	"meal-planner/web"
	"meal-planner/web/views"
)

func main() {
	ctx := core.InitApp()
	defer ctx.Close()

	//core.InitDb(ctx)
	//irc.InitDb(ctx)
	//shows.InitDb(ctx)

	e := echo.New()
	e.Renderer = web.NewTemplate(ctx)
	e.Debug = true
	//e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(core.CreateCtx(ctx))

	root := e.Group(ctx.AppConfig.BasePath)

	root.Static("/assets", "./web/assets")

	root.GET("/", views.Index)
	root.GET("/offset/:offset", views.Offset)

	// Listen and server on 0.0.0.0:8080
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", ctx.AppConfig.Port)))
}
