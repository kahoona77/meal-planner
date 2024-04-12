package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"meal-planner/core"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const templatesDir = "web/tmpl"

type HtmlRenderer struct {
	templates    map[string]*template.Template
	basePath     string
	isDev        bool
	manifest     Manifest
	devServerUrl string
}

func NewRenderer(basePath string, manifest Manifest, isDev bool) *HtmlRenderer {
	renderer := HtmlRenderer{
		templates: map[string]*template.Template{},
		basePath:  basePath,
		manifest:  manifest,
		isDev:     isDev,
	}

	if isDev {
		renderer.devServerUrl = "http://localhost:5173"
	}

	renderer.loadTemplates()
	return &renderer
}

func CreateRenderer(ctx *core.Ctx) (core.HtmlRenderer, error) {
	manifest := EmptyManifest()
	f := os.DirFS("./")

	if !ctx.Config().IsDev {
		var err error
		manifest, err = ParseManifest("web/assets/dist/manifest.json", f)
		if err != nil {
			return nil, fmt.Errorf("could not prase manifest-file: %v", err)
		}
	}

	return NewRenderer(ctx.Config().BasePath, manifest, ctx.Config().IsDev), nil
}

func (t *HtmlRenderer) loadTemplates() {
	ext := ".html"
	layout := templatesDir + "/base" + ext
	_, err := os.Stat(layout)
	if err != nil {
		logrus.Panicf("cannot find %s", layout)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"basePath": func() string {
			return fmt.Sprintf("%s/", t.basePath)
		},
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"isToday": func(date time.Time) bool {
			return IsToday(date)
		},
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
		"fileUrl": func(id sql.NullInt64) string {
			if !id.Valid {
				return ""
			}
			return fmt.Sprintf("files/%d", id.Int64)
		},
		"assetUrl": func(asset string) string {
			if !t.isDev {
				asset = t.manifest.File(asset)
			}

			return fmt.Sprintf("%s%s/assets/%s", t.devServerUrl, t.basePath, asset)
		},
		"publicUrl": func(asset string) string {
			if !t.isDev {
				asset = t.manifest.File(asset)
			}

			serverUrl := t.devServerUrl
			if strings.Contains(asset, ".svg") {
				serverUrl = ""
			}

			return fmt.Sprintf("%s%s/%s", serverUrl, t.basePath, asset)
		},
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"formatWeekday": func(weekday interface{}) string {
			if w, ok := weekday.(int); ok {
				return formatWeekday(w)
			}
			if w, ok := weekday.(time.Weekday); ok {
				return formatWeekday(int(w))
			}

			return ""
		},
	}

	views, _ := filepath.Glob(templatesDir + "**/*" + ext)

	// first find all partials
	partials := make([]string, 0)
	for _, view := range views {
		_, file := filepath.Split(view)
		if strings.HasPrefix(file, "_") {
			partials = append(partials, view)
		}
	}

	for _, view := range views {
		_, file := filepath.Split(view)
		//dir = strings.Replace(dir, templatesDir, "", 1)
		//file = strings.TrimSuffix(file, ext)
		renderName := file

		tmplFiles := append([]string{layout, view}, partials...)
		tmpl := template.Must(template.New(filepath.Base(layout)).Funcs(funcMap).ParseFiles(tmplFiles...))
		t.Add(renderName, tmpl)
	}
}

func (t *HtmlRenderer) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	t.templates[name] = tmpl
}

func (t *HtmlRenderer) Render(w io.Writer, name string, data core.TemplateData, ctx *core.WebContext) error {
	t.loadTemplates()
	if _, ok := t.templates[name]; ok == false {
		return fmt.Errorf("no such view. (%s)", name)
	}

	data["isDev"] = t.isDev
	data["basePath"] = t.basePath
	data["manifest"] = t.manifest

	return t.templates[name].Execute(w, data)
}

func IsToday(date time.Time) bool { //get monday 00:00:00
	now := time.Now()
	return now.Year() == date.Year() && now.Month() == date.Month() && now.Day() == date.Day()
}
