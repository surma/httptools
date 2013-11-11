package httptools

import (
	"net/http"
	"path"
	"strings"
)

// DiscardPathElements discards n elements from the request path.
// It's most useful in a handler list. The original request path
// can be found in the VarsResponseWriter under "OrigPath".
func DiscardPathElements(n int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vrw, ok := w.(VarsResponseWriter)
		if !ok {
			vrw = newOurResponseWriter(w)
		}
		if _, ok := vrw.Vars()["OrigPath"]; !ok {
			vrw.Vars()["OrigPath"] = r.URL.Path
		}

		elems := strings.Split(r.URL.Path, "/")
		if n >= len(elems) {
			r.URL.Path = "/"
			return
		}
		r.URL.Path = "/" + path.Join(elems[n:]...)
	})
}
