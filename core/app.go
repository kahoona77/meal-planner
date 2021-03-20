package core

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func InitApp() *App {
	formatter := &logrus.TextFormatter{}
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	conf := LoadConfiguration()

	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", conf.DbFile)
	if err != nil {
		panic(err)
	}

	ctx := &Ctx{config: &conf, db: db}

	e := echo.New()
	e.Renderer = NewTemplate(ctx)
	e.Debug = true
	//e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(CreateCtx(ctx))

	return &App{Echo: e, Ctx: ctx}
}

type App struct {
	*echo.Echo
	Ctx *Ctx
}

type HandlerFunc func(*WebContext) error

func (f HandlerFunc) Handle(ctx *WebContext) error {
	return f(ctx)
}
func wrapHandler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*WebContext)
		return h.Handle(ctx)
	}
}

func (a *App) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return a.Add(http.MethodGet, path, wrapHandler(h), m...)
}

func (a *App) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return a.Add(http.MethodPost, path, wrapHandler(h), m...)
}

func (a *App) PUT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return a.Add(http.MethodPut, path, wrapHandler(h), m...)
}

func (a *App) DELETE(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return a.Add(http.MethodPut, path, wrapHandler(h), m...)
}

func (a *App) Group(prefix string, m ...echo.MiddlewareFunc) *Group {
	g := a.Echo.Group(prefix, m...)
	return &Group{Group: g}
}

type Group struct {
	*echo.Group
}

func (g *Group) GET(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return g.Add(http.MethodGet, path, wrapHandler(h), m...)
}

func (g *Group) POST(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return g.Add(http.MethodPost, path, wrapHandler(h), m...)
}

func (g *Group) PUT(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return g.Add(http.MethodPut, path, wrapHandler(h), m...)
}

func (g *Group) DELETE(path string, h HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return g.Add(http.MethodPut, path, wrapHandler(h), m...)
}
