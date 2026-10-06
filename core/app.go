package core

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

func InitApp(rendererFactory CreateRendererFunc, migrationsFs fs.FS) *App {
	conf := LoadConfiguration()

	slog.SetDefault(slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:      conf.LogLevel,
		TimeFormat: time.DateTime,
		NoColor:    !isTerminal(os.Stdout),
	})))

	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", conf.DbFile)
	if err != nil {
		panic(err)
	}

	if err := runMigrations(db.DB, migrationsFs); err != nil {
		panic(err)
	}

	ctx := &Ctx{config: &conf, db: db}
	renderer, err := rendererFactory(ctx)
	if err != nil {
		panic(err)
	}

	return &App{mux: http.NewServeMux(), Ctx: ctx, renderer: renderer}
}

type App struct {
	mux      *http.ServeMux
	renderer HtmlRenderer
	Ctx      *Ctx
}

// Start listens on addr and serves all registered routes.
func (a *App) Start(addr string) error {
	slog.Info("http server started", "addr", addr)
	return http.ListenAndServe(addr, logRequests(a.mux))
}

func (a *App) Group(prefix string) *Group {
	prefix = strings.TrimSuffix(prefix, "/")
	if prefix != "" {
		// the group root without trailing slash, e.g. /meal-planner -> /meal-planner/
		a.mux.Handle("GET "+prefix, http.RedirectHandler(prefix+"/", http.StatusMovedPermanently))
	}
	return &Group{app: a, prefix: prefix}
}

type HandlerFunc func(*WebContext) error

type Group struct {
	app    *App
	prefix string
}

func (g *Group) pattern(method string, path string) string {
	if path == "/" {
		// match only the root itself, not everything below it
		path = "/{$}"
	}
	return method + " " + g.prefix + path
}

func (g *Group) handle(method string, path string, h HandlerFunc) {
	g.app.mux.HandleFunc(g.pattern(method, path), func(w http.ResponseWriter, r *http.Request) {
		ctx := &WebContext{Ctx: g.app.Ctx, renderer: g.app.renderer, w: w, r: r}
		if err := h(ctx); err != nil {
			handleError(w, r, err)
		}
	})
}

func (g *Group) GET(path string, h HandlerFunc) {
	g.handle(http.MethodGet, path, h)
}

func (g *Group) POST(path string, h HandlerFunc) {
	g.handle(http.MethodPost, path, h)
}

func (g *Group) PUT(path string, h HandlerFunc) {
	g.handle(http.MethodPut, path, h)
}

func (g *Group) DELETE(path string, h HandlerFunc) {
	g.handle(http.MethodDelete, path, h)
}

// Static serves the files in dir under path, without directory listings.
func (g *Group) Static(path string, dir string) {
	prefix := g.prefix + strings.TrimSuffix(path, "/") + "/"
	fileServer := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
	g.app.mux.Handle("GET "+prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	}))
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	slog.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info(r.Method+" "+r.URL.Path, "status", rec.status, "duration", time.Since(start).Round(time.Microsecond))
	})
}

// isTerminal reports whether f is an interactive terminal, so colors are only used there.
func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
