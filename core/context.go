package core

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
)

type Context interface {
	Db() *sqlx.DB
	Config() *AppConfig
}

func CreateCtx(ctx *Ctx, renderer HtmlRenderer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &WebContext{Context: c, Ctx: ctx, renderer: renderer}
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
	renderer HtmlRenderer
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

func (ctx *WebContext) RenderTemplate(code int, name string, data TemplateData) (err error) {
	if ctx.renderer == nil {
		return echo.ErrRendererNotRegistered
	}
	buf := new(bytes.Buffer)
	if err = ctx.renderer.Render(buf, name, data, ctx); err != nil {
		return
	}
	return ctx.HTMLBlob(code, buf.Bytes())
}

type HtmlRenderer interface {
	Render(w io.Writer, name string, data TemplateData, ctx *WebContext) error
}

type CreateRendererFunc func(ctx *Ctx) (HtmlRenderer, error)

type TemplateData map[string]any
