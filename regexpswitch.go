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

func NewRegexpSwitch(routes map[string]http.Handler) http.Handler {
	rs := regexpSwitch{}
	for re, h := range routes {
		rs[MustRegexp("^"+re+"$")] = h
	}
	return rs
}

func MustRegexp(re string) *regexp.Regexp {
	r, err := regexp.CompilePOSIX(re)
	if err != nil {
		panic(err)
	}
	return r
}
