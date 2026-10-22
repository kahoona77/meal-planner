package web

import (
	"context"
	"net/http"
	"strings"
)

// page holds what every rendered page needs, independent of the handler.
type page struct {
	basePath string
	assets   *Assets
	path     string
}

type pageKey struct{}

// PageMiddleware makes the base path, the asset URLs and the current path available to
// the templates (see BasePath, Asset and Path).
func PageMiddleware(basePath string, assets *Assets) func(http.Handler) http.Handler {
	p := &page{basePath: basePath, assets: assets}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rp := *p
			rp.path = strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, basePath), "/")
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pageKey{}, &rp)))
		})
	}
}

func pageFrom(ctx context.Context) *page {
	if p, ok := ctx.Value(pageKey{}).(*page); ok {
		return p
	}
	panic("web.PageMiddleware is not installed")
}

// BasePath returns the base path with a trailing slash, for <base href>.
func BasePath(ctx context.Context) string {
	return pageFrom(ctx).basePath + "/"
}

// Asset returns the versioned URL of a static file, e.g. Asset(ctx, "img/solid.svg#plus").
func Asset(ctx context.Context, path string) string {
	return pageFrom(ctx).assets.URL(path)
}

// Path returns the path of the current request relative to the base path, without a
// leading slash, e.g. "meals/1".
func Path(ctx context.Context) string {
	return pageFrom(ctx).path
}
