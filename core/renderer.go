package core

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const templatesDir = "web/tmpl"

type Template struct {
	templates map[string]*template.Template
	basePath  string
}

func NewTemplate(ctx Context) *Template {
	ins := Template{
		templates: map[string]*template.Template{},
		basePath:  ctx.Config().BasePath,
	}

	ins.loadTemplates()
	return &ins
}

func (t *Template) loadTemplates() {
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

func (t *Template) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	t.templates[name] = tmpl
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	t.loadTemplates()
	if _, ok := t.templates[name]; ok == false {
		return fmt.Errorf("no such view. (%s)", name)
	}
	return t.templates[name].Execute(w, data)
}

func IsToday(date time.Time) bool { //get monday 00:00:00
	now := time.Now()
	return now.Year() == date.Year() && now.Month() == date.Month() && now.Day() == date.Day()
}
