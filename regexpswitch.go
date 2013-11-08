package httptools

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
)

type regexpRule struct {
	*regexp.Regexp
	http.Handler
}

// RegexpRule represents a single rule in a RegexpSwitch.
type RegexpRule interface {
	// Same method provided by regexp.Regexp.
	// The returned array will be saved to the VarsResponseWriter.
	FindStringSubmatch(s string) []string
	http.Handler
}

// RegexpSwitch is a slice of RegexpRules. They will be checked
// in the order they have been provided. If a rule matches
// (i.e. Regexp.Rule.FindStringSubmatch return value is non-nil), the
// Handler will be called and the slice traversal is stopped.
type RegexpSwitch []RegexpRule

func (rs RegexpSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orw, ok := w.(*ourResponseWriter)
	if !ok {
		orw = newOurResponseWriter(w)
	}

	for _, rule := range rs {
		if m := rule.FindStringSubmatch(r.URL.Path); m != nil {
			for i := 1; i < len(m); i++ {
				orw.Vars()[fmt.Sprintf("%d", i)] = m[i]
			}
			rule.ServeHTTP(orw, r)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

type regexpSwitch []regexpRule

func (rs regexpSwitch) Len() int {
	return len(rs)
}

func (rs regexpSwitch) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs regexpSwitch) Less(i, j int) bool {
	return len(rs[i].String()) < len(rs[j].String())
}

// A regexp switch takes a map of regexp strings and handlers.
// A regexp is considered a match if it matches the whole string.
// Longer patterns take precedence over shorter ones.
func NewRegexpSwitch(routes map[string]http.Handler) RegexpSwitch {
	rs := make(regexpSwitch, 0, len(routes))
	for re, h := range routes {
		rs = append(rs, regexpRule{
			Regexp:  regexp.MustCompilePOSIX("^" + re + "$"),
			Handler: h,
		})
	}
	sort.Sort(sort.Reverse(rs))

	nrs := make(RegexpSwitch, 0, len(routes))
	for _, rr := range rs {
		nrs = append(nrs, RegexpRule(rr))
	}
	return nrs
}
