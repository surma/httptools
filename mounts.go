package httptools

import (
	"net/http"
	"net/url"
	"strings"
)

// Mounts is a list of handlers which are mounted at the given path.
// Mounting is a simple path prefix-based routing. The prefix will be
// stripped from the request before being passed to the associated handler.
type Mounts map[string]http.Handler

func (m Mounts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for path, h := range m {
		path = strings.TrimRight(path, "/")
		if strings.HasPrefix(r.URL.Path, path) {
			// FIXME: This is a shallow copy
			nr := *r
			nr.URL = &url.URL{}
			*nr.URL = *r.URL

			nr.URL.Path = strings.TrimPrefix(nr.URL.Path, path)
			if h != nil {
				h.ServeHTTP(w, &nr)
			}
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}
