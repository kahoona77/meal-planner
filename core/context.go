package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
)

// maxMemory is the part of a multipart form kept in memory, the rest goes to temp files.
const maxMemory = 32 << 20

type Context interface {
	Db() *sqlx.DB
	Config() *AppConfig
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
		slog.Error("error closing database", "err", err)
	}
}

// WebContext is passed to every handler and wraps the current request and response.
type WebContext struct {
	*Ctx
	renderer HtmlRenderer
	w        http.ResponseWriter
	r        *http.Request
}

func (ctx *WebContext) Request() *http.Request {
	return ctx.r
}

func (ctx *WebContext) Header() http.Header {
	return ctx.w.Header()
}

// Param returns the path parameter name, e.g. {id}, or "" if the route has none.
func (ctx *WebContext) Param(name string) string {
	return ctx.r.PathValue(name)
}

func (ctx *WebContext) ParamAsInt(name string) int {
	param := ctx.Param(name)
	p, _ := strconv.Atoi(param)
	return p
}

func (ctx *WebContext) FormValue(name string) string {
	return ctx.r.FormValue(name)
}

func (ctx *WebContext) FormFile(name string) (*multipart.FileHeader, error) {
	if err := ctx.r.ParseMultipartForm(maxMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return nil, err
	}
	f, header, err := ctx.r.FormFile(name)
	if err != nil {
		return nil, err
	}
	return header, f.Close()
}

// Redirect redirects to path, which is relative to the configured base path.
func (ctx *WebContext) Redirect(code int, path string) error {
	http.Redirect(ctx.w, ctx.r, fmt.Sprintf("%s%s", ctx.config.BasePath, path), code)
	return nil
}

func (ctx *WebContext) Blob(code int, contentType string, data []byte) error {
	ctx.w.Header().Set("Content-Type", contentType)
	ctx.w.WriteHeader(code)
	_, err := ctx.w.Write(data)
	return err
}

func (ctx *WebContext) RenderTemplate(code int, name string, data TemplateData) (err error) {
	if ctx.renderer == nil {
		return errors.New("no renderer registered")
	}
	buf := new(bytes.Buffer)
	if err = ctx.renderer.Render(buf, name, data, ctx); err != nil {
		return
	}
	return ctx.Blob(code, "text/html; charset=UTF-8", buf.Bytes())
}

type HtmlRenderer interface {
	Render(w io.Writer, name string, data TemplateData, ctx *WebContext) error
}

type CreateRendererFunc func(ctx *Ctx) (HtmlRenderer, error)

type TemplateData map[string]any
