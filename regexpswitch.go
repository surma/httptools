package httptools

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
)

type regexpRule struct {
	re *regexp.Regexp
	h  http.Handler
}

type regexpSwitch []regexpRule

func (rs regexpSwitch) Len() int {
	return len(rs)
}

func (rs regexpSwitch) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs regexpSwitch) Less(i, j int) bool {
	return len(rs[i].re.String()) < len(rs[j].re.String())
}

func (rs regexpSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orw, ok := w.(*ourResponseWriter)
	if !ok {
		orw = newOurResponseWriter(w)
	}

	for _, rule := range rs {
		if m := rule.re.FindStringSubmatch(r.URL.Path); m != nil {
			for i := 1; i < len(m); i++ {
				orw.Vars()[fmt.Sprintf("%d", i)] = m[i]
			}
			rule.h.ServeHTTP(orw, r)
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
		rs = append(rs, regexpRule{
			re: mustRegexp("^" + re + "$"),
			h:  h,
		})
	}
	sort.Sort(sort.Reverse(rs))
	return rs
}

func mustRegexp(re string) *regexp.Regexp {
	r, err := regexp.CompilePOSIX(re)
	if err != nil {
		panic(err)
	}
	return r
}
