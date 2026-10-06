package core

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"
)

func InitApp(rendererFactory CreateRendererFunc, migrationsFs fs.FS) *App {
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
	logrus.Infof("http server started on %s", addr)
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
	logrus.Errorf("%s %s: %v", r.Method, r.URL.Path, err)
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
		logrus.Infof("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
	})
}
