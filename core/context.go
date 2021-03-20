package core

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Context interface {
	Db() *sqlx.DB
	Config() *AppConfig
}

func CreateCtx(ctx *Ctx) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &WebContext{Context: c, Ctx: ctx}
			return next(cc)
		}
	}
}

type Ctx struct {
	db     *sqlx.DB
	config *AppConfig
}

func (ctx *Ctx) Db() *sqlx.DB {
	return ctx.db
}

func (ctx *Ctx) Config() *AppConfig {
	return ctx.config
}

func (ctx *Ctx) Close() {
	if err := ctx.db.Close(); err != nil {
		logrus.Errorf("error closing database: %v", err)
	}
}

type WebContext struct {
	echo.Context
	*Ctx
}

func (ctx *WebContext) ParamAsInt(name string) int {
	param := ctx.Param(name)
	p, _ := strconv.Atoi(param)
	return p
}

func (ctx *WebContext) Redirect(code int, name string) error {
	return ctx.Context.Redirect(code, fmt.Sprintf("%s%s", ctx.config.BasePath, name))
}
