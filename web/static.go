package web

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// Assets builds URLs for the static files. Each URL carries a hash of the file
// content (?v=...), so browsers can cache it forever and still get new versions.
type Assets struct {
	fsys     fs.FS
	basePath string

	mu       sync.Mutex
	versions map[string]assetVersion
}

type assetVersion struct {
	modTime time.Time
	size    int64
	hash    string
}

func NewAssets(fsys fs.FS, basePath string) *Assets {
	return &Assets{fsys: fsys, basePath: basePath, versions: map[string]assetVersion{}}
}

// URL returns the URL of the static file at path, e.g. "img/solid.svg#plus".
// A #fragment is kept, a missing file gets a URL without version.
func (a *Assets) URL(path string) string {
	file, fragment, _ := strings.Cut(path, "#")
	url := fmt.Sprintf("%s/assets/%s", a.basePath, file)
	if v := a.version(file); v != "" {
		url += "?v=" + v
	}
	if fragment != "" {
		url += "#" + fragment
	}
	return url
}

// version returns the content hash of file. It is cached until the file changes,
// so a rebuilt CSS file in development gets a new URL right away.
func (a *Assets) version(file string) string {
	info, err := fs.Stat(a.fsys, file)
	if err != nil {
		slog.Warn("static file not found", "path", file)
		return ""
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if v, ok := a.versions[file]; ok && v.modTime.Equal(info.ModTime()) && v.size == info.Size() {
		return v.hash
	}

	data, err := fs.ReadFile(a.fsys, file)
	if err != nil {
		slog.Warn("could not read static file", "path", file, "err", err)
		return ""
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])[:8]
	a.versions[file] = assetVersion{modTime: info.ModTime(), size: info.Size(), hash: hash}
	return hash
}
