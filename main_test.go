package main

import (
	"bytes"
	"io"
	"meal-planner/core"
	"meal-planner/internal/testdb"
	"meal-planner/planner"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const basePath = "/meal-planner"

// testApp is the complete app (all routes and middlewares) on an in-memory database.
type testApp struct {
	t      *testing.T
	db     *sqlx.DB
	server *httptest.Server
	client *http.Client
}

func newTestApp(t *testing.T) *testApp {
	db := testdb.New(t)
	db.MustExec(`INSERT INTO tags (id, name, color) VALUES (1, 'Suppe', '#61187c'), (2, 'Nudeln', '#fbe3a7')`)
	db.MustExec(`INSERT INTO meals (id, name, description) VALUES
		(1, 'Tomatensuppe', '<div>heiß</div>'), (2, 'Spaghetti', ''), (3, 'Lasagne', ''), (4, 'Pfannkuchen', ''),
		(5, 'Burger', ''), (6, 'Gratin', ''), (7, 'Eintopf', '')`)
	db.MustExec(`INSERT INTO meal_tags (meal_id, tag_id) VALUES (1, 1), (2, 2), (3, 2)`)

	app := core.NewApp(core.AppConfig{BasePath: basePath}, db)
	setupRoutes(app)
	server := httptest.NewServer(app.Handler())
	t.Cleanup(server.Close)

	// don't follow redirects, the tests check them
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &testApp{t: t, db: db, server: server, client: client}
}

func (a *testApp) do(req *http.Request) (*http.Response, string) {
	a.t.Helper()
	res, err := a.client.Do(req)
	require.NoError(a.t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(a.t, err)
	return res, string(body)
}

func (a *testApp) get(path string) (*http.Response, string) {
	a.t.Helper()
	req, err := http.NewRequest(http.MethodGet, a.server.URL+path, nil)
	require.NoError(a.t, err)
	return a.do(req)
}

func (a *testApp) postForm(path string, form url.Values) *http.Response {
	a.t.Helper()
	req, err := http.NewRequest(http.MethodPost, a.server.URL+path, strings.NewReader(form.Encode()))
	require.NoError(a.t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, _ := a.do(req)
	return res
}

// postMultipart posts fields like the meal form does; image is uploaded as file "image" if not nil.
func (a *testApp) postMultipart(path string, fields url.Values, image []byte) *http.Response {
	a.t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for name, values := range fields {
		for _, v := range values {
			require.NoError(a.t, w.WriteField(name, v))
		}
	}
	if image != nil {
		part, err := w.CreateFormFile("image", "meal.png")
		require.NoError(a.t, err)
		_, err = part.Write(image)
		require.NoError(a.t, err)
	}
	require.NoError(a.t, w.Close())

	req, err := http.NewRequest(http.MethodPost, a.server.URL+path, body)
	require.NoError(a.t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, _ := a.do(req)
	return res
}

func (a *testApp) count(query string, args ...any) int {
	a.t.Helper()
	var n int
	require.NoError(a.t, a.db.Get(&n, query, args...))
	return n
}

func assertRedirect(t *testing.T, res *http.Response, location string) {
	t.Helper()
	assert.Equal(t, http.StatusFound, res.StatusCode)
	assert.Equal(t, location, res.Header.Get("Location"))
}

func TestPages(t *testing.T) {
	app := newTestApp(t)
	app.db.MustExec(`INSERT INTO meals_of_day (id, date, meal_id) VALUES ('2026_10_05', '2026-10-05', 1)`)

	tests := []struct {
		path     string
		contains string
	}{
		{"/", "Create new"},
		{"/", `href="planner/` + time.Now().Format(planner.MealOfDayDateFormat) + `/select"`}, // nothing planned yet
		{"/offset/-1", "Create new"},
		{"/planner/2026_10_05", "<div>heiß</div>"},
		{"/planner/2026_10_05/select", `value="1" checked`},
		{"/meals", "Tomatensuppe"},
		{"/meals/new", `action="meals"`},
		{"/meals/1", `value="Tomatensuppe"`},
		{"/tags", "Nudeln"},
		{"/wizard/0", "Generate"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			res, body := app.get(basePath + tt.path)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "text/html; charset=UTF-8", res.Header.Get("Content-Type"))
			assert.Contains(t, body, tt.contains)
			assert.Contains(t, body, `<base href="/meal-planner/">`)
		})
	}
}

func TestWeekShowsNextMeal(t *testing.T) {
	app := newTestApp(t)
	planMeal := func(date time.Time, mealId int) {
		app.db.MustExec(`INSERT OR REPLACE INTO meals_of_day (id, date, meal_id) VALUES (?, ?, ?)`,
			date.Format(planner.MealOfDayDateFormat), date.Format(time.DateOnly), mealId)
	}

	_, body := app.get(basePath + "/")
	assert.NotContains(t, body, "<div>heiß</div>")

	// the next planned day is shown in full, with tags and description
	tomorrow := time.Now().AddDate(0, 0, 1)
	if planner.GetStartWeek(tomorrow).After(time.Now()) {
		t.Skip("tomorrow is in the next week")
	}
	planMeal(tomorrow, 1)
	_, body = app.get(basePath + "/")
	assert.Contains(t, body, "Als Nächstes")
	assert.Contains(t, body, "<div>heiß</div>")
	assert.Contains(t, body, `--tag-color: #61187c`)

	// today comes first
	planMeal(time.Now(), 2)
	_, body = app.get(basePath + "/")
	assert.Contains(t, body, "Heute")
	assert.NotContains(t, body, "<div>heiß</div>")
}

func TestActiveNavLink(t *testing.T) {
	app := newTestApp(t)

	tests := map[string]string{"/": "", "/planner/2026_10_05": "", "/meals": "meals", "/meals/1": "meals", "/tags": "tags"}
	for path, href := range tests {
		_, body := app.get(basePath + path)
		assert.Contains(t, body, `<a href="`+href+`" aria-current="page"`, path)
		assert.Equal(t, 1, strings.Count(body, `aria-current="page"`), path)
	}
}

func TestNotFound(t *testing.T) {
	app := newTestApp(t)

	for _, path := range []string{"/meals/999", "/files/999", "/unknown", "/assets/missing.css", "/assets/img/"} {
		res, _ := app.get(basePath + path)
		assert.Equal(t, http.StatusNotFound, res.StatusCode, path)
	}
}

func TestBasePathRedirect(t *testing.T) {
	app := newTestApp(t)

	res, _ := app.get(basePath)

	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)
	assert.Equal(t, basePath+"/", res.Header.Get("Location"))
}

func TestStaticAssets(t *testing.T) {
	app := newTestApp(t)

	_, page := app.get(basePath + "/")
	versioned := regexp.MustCompile(`/meal-planner/assets/htmx\.min\.js\?v=[0-9a-f]{8}`).FindString(page)
	require.NotEmpty(t, versioned)

	res, body := app.get(versioned)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "public, max-age=31536000, immutable", res.Header.Get("Cache-Control"))
	assert.NotEmpty(t, body)

	res, _ = app.get(basePath + "/assets/htmx.min.js")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "no-cache", res.Header.Get("Cache-Control"))
}

func TestCreateMeal(t *testing.T) {
	app := newTestApp(t)
	image := []byte("\x89PNG fake image")

	res := app.postMultipart(basePath+"/meals", url.Values{
		"name":        {"Kürbissuppe"},
		"description": {"<div>lecker</div>"},
		"tags":        {"1", "2"},
	}, image)

	assertRedirect(t, res, basePath+"/meals")
	var mealId, imageId int64
	require.NoError(t, app.db.QueryRow(`SELECT id, image_file_id FROM meals WHERE name = 'Kürbissuppe'`).Scan(&mealId, &imageId))
	assert.Equal(t, 2, app.count(`SELECT count(*) FROM meal_tags WHERE meal_id = $1`, mealId))

	fileRes, data := app.get(basePath + "/files/" + strconv.FormatInt(imageId, 10))
	assert.Equal(t, http.StatusOK, fileRes.StatusCode)
	assert.Equal(t, string(image), data)
	assert.Equal(t, "max-age=432000", fileRes.Header.Get("Cache-Control"))
}

func TestEditMeal(t *testing.T) {
	app := newTestApp(t)

	// no image: the meal keeps its (missing) image, tags are replaced
	res := app.postMultipart(basePath+"/meals/2", url.Values{"name": {"Spaghetti Bolognese"}, "tags": {"1"}}, nil)

	assertRedirect(t, res, basePath+"/meals")
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM meals WHERE id = 2 AND name = 'Spaghetti Bolognese'`))
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM meal_tags WHERE meal_id = 2 AND tag_id = 1`))
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM meal_tags WHERE meal_id = 2`))
}

func TestDeleteMeal(t *testing.T) {
	app := newTestApp(t)
	app.postMultipart(basePath+"/meals/1", url.Values{"name": {"Tomatensuppe"}, "tags": {"1"}}, []byte("img"))
	require.Equal(t, 1, app.count(`SELECT count(*) FROM files`))

	res := app.postForm(basePath+"/meals/1/delete", nil)

	assertRedirect(t, res, basePath+"/meals")
	assert.Equal(t, 0, app.count(`SELECT count(*) FROM meals WHERE id = 1`))
	assert.Equal(t, 0, app.count(`SELECT count(*) FROM meal_tags WHERE meal_id = 1`))
	assert.Equal(t, 0, app.count(`SELECT count(*) FROM files`))
}

func TestTags(t *testing.T) {
	app := newTestApp(t)

	res := app.postForm(basePath+"/tags", url.Values{"name": {"Asiatisch"}, "color": {"#e36887"}})
	assertRedirect(t, res, basePath+"/tags")
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM tags WHERE name = 'Asiatisch' AND color = '#e36887'`))

	res = app.postForm(basePath+"/tags/1", url.Values{"name": {"Suppen"}, "color": {"#000000"}})
	assertRedirect(t, res, basePath+"/tags")
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM tags WHERE id = 1 AND name = 'Suppen' AND color = '#000000'`))
}

func TestSelectMealOfDay(t *testing.T) {
	app := newTestApp(t)

	res := app.postForm(basePath+"/planner/2026_10_05/select", url.Values{"selected": {"3"}})
	assertRedirect(t, res, basePath+"/")
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM meals_of_day WHERE id = '2026_10_05' AND meal_id = 3`))

	// "- leer -" clears the day
	app.postForm(basePath+"/planner/2026_10_05/select", url.Values{"selected": {""}})
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM meals_of_day WHERE id = '2026_10_05' AND meal_id IS NULL`))
}

func TestWizard(t *testing.T) {
	app := newTestApp(t)

	// the wizard only picks meals with tags, and it fills the days in order, so every day gets
	// tags here: Monday the only soup, all other days one of the six pasta meals
	app.db.MustExec(`INSERT INTO meal_tags (meal_id, tag_id) VALUES (4, 2), (5, 2), (6, 2), (7, 2)`)
	form := url.Values{"tags_1": {"1"}}
	for _, weekday := range []string{"0", "2", "3", "4", "5", "6"} {
		form.Set("tags_"+weekday, "2")
	}
	res := app.postForm(basePath+"/wizard/0", form)

	assertRedirect(t, res, basePath+"/offset/0")
	assert.Equal(t, 1, app.count(`SELECT count(*) FROM weekday_tags WHERE weekday = 1 AND tag_id = 1`))

	var mealIds []int64
	require.NoError(t, app.db.Select(&mealIds, `SELECT meal_id FROM meals_of_day ORDER BY date`))
	assert.Len(t, mealIds, 7, "one meal per day")
	assert.Len(t, uniq(mealIds), 7, "no meal twice in a week")

	var soupDay string
	require.NoError(t, app.db.Get(&soupDay, `SELECT id FROM meals_of_day WHERE meal_id = 1`))
	date, err := time.Parse(planner.MealOfDayDateFormat, soupDay)
	require.NoError(t, err)
	assert.Equal(t, time.Monday, date.Weekday(), "the only soup is on Monday")

	_, page := app.get(basePath + "/wizard/0")
	assert.Regexp(t, `name="tags_1" value="1" checked`, page, "the wizard remembers the tags")
}

func uniq(ids []int64) map[int64]bool {
	set := map[int64]bool{}
	for _, id := range ids {
		set[id] = true
	}
	return set
}
