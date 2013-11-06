package httptools

import (
	"fmt"
	"net/http"
	"regexp"
)

type RegexpSwitch map[*regexp.Regexp]http.Handler

func (rs RegexpSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func MustRegexp(re string) *regexp.Regexp {
	r, err := regexp.CompilePOSIX(re)
	if err != nil {
		panic(err)
	}
	return r
}
