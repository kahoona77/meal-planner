package web

import (
	"context"
	"net/http"
)

// page holds what every rendered page needs, independent of the handler.
type page struct {
	basePath string
	assets   *Assets
}

type pageKey struct{}

// PageMiddleware makes the base path and the asset URLs available to the templates
// (see BasePath and Asset).
func PageMiddleware(basePath string, assets *Assets) func(http.Handler) http.Handler {
	p := &page{basePath: basePath, assets: assets}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pageKey{}, p)))
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
