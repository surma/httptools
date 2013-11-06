package httptools

import (
	"fmt"
	"net/http"
	"regexp"
)

type regexpSwitch map[*regexp.Regexp]http.Handler

func (rs regexpSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orw, ok := w.(*ourResponseWriter)
	if !ok {
		orw = newOurResponseWriter(w)
	}

	for re, h := range rs {
		if m := re.FindStringSubmatch(r.URL.Path); m != nil {
			for i := 1; i < len(m); i++ {
				orw.Vars()[fmt.Sprintf("%d", i)] = m[i]
			}
			h.ServeHTTP(orw, r)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

// A regexp switch takes a map of regexp strings and handlers.
// If a request path matches a regexp, the corresponding handler is
// executed. Submatches will be put inside a VarsResponseWriter with the
// keys "1", "2", ...
func NewRegexpSwitch(routes map[string]http.Handler) http.Handler {
	rs := regexpSwitch{}
	for re, h := range routes {
		rs[mustRegexp("^"+re+"$")] = h
	}
	return rs
}

func mustRegexp(re string) *regexp.Regexp {
	r, err := regexp.CompilePOSIX(re)
	if err != nil {
		panic(err)
	}
	return r
}
