package core

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"strconv"
)

func CreateCtx(ctx *Ctx) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &WebContext{Context: c, Ctx: ctx}
			return next(cc)
		}
	}
}

type Ctx struct {
	Db         *sqlx.DB
	AppConfig  *AppConfig
	IrcService IrcService
}

func (ctx *Ctx) Copy() *Ctx {
	return &Ctx{Db: ctx.Db, IrcService: ctx.IrcService}
}

func (ctx *Ctx) Close() {
	if err := ctx.Db.Close(); err != nil {
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
	return ctx.Context.Redirect(code, fmt.Sprintf("%s%s", ctx.AppConfig.BasePath, name))
}
