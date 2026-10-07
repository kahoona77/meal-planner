package web

import (
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAssets_URL(t *testing.T) {
	fsys := fstest.MapFS{
		"main.css":      {Data: []byte("body{}"), ModTime: time.Unix(1, 0)},
		"img/solid.svg": {Data: []byte("<svg/>"), ModTime: time.Unix(1, 0)},
	}
	assets := NewAssets(fsys, "/meal-planner")

	css := assets.URL("main.css")
	assert.Regexp(t, `^/meal-planner/assets/main\.css\?v=[0-9a-f]{8}$`, css)

	svg := assets.URL("img/solid.svg#plus")
	assert.Regexp(t, `^/meal-planner/assets/img/solid\.svg\?v=[0-9a-f]{8}#plus$`, svg)

	assert.Equal(t, "/meal-planner/assets/missing.js", assets.URL("missing.js"))
}

func TestAssets_URL_changesWithContent(t *testing.T) {
	fsys := fstest.MapFS{"main.css": {Data: []byte("body{}"), ModTime: time.Unix(1, 0)}}
	assets := NewAssets(fsys, "")

	before := assets.URL("main.css")
	assert.Equal(t, before, assets.URL("main.css"))

	fsys["main.css"] = &fstest.MapFile{Data: []byte("body{color:red}"), ModTime: time.Unix(2, 0)}
	after := assets.URL("main.css")

	assert.NotEqual(t, before, after)
	assert.True(t, strings.HasPrefix(after, "/assets/main.css?v="))
}
