package httptools

import (
	"net/http"
	"regexp"
	"sort"
)

// RegexpPathRule is a PathRule using a regular expression.
type RegexpPathRule struct {
	*regexp.Regexp
	http.Handler
}

func (rpr RegexpPathRule) Match(s string) ([]string, bool) {
	submatches := rpr.FindStringSubmatch(s)
	if submatches == nil || len(submatches) < 1 {
		return nil, false
	}
	return submatches[1:], true
}

func (rpr RegexpPathRule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rpr.Handler != nil {
		rpr.Handler.ServeHTTP(w, r)
	}
}

type regexpSwitch []RegexpPathRule

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
// A regexp is considered a match if it matches the string from the
// beginning (but not necessarily to the end).
// Longer patterns take precedence over shorter ones.
func NewRegexpSwitch(routes map[string]http.Handler) PathSwitch {
	rs := make(regexpSwitch, 0, len(routes))
	for re, h := range routes {
		rs = append(rs, RegexpPathRule{
			Regexp:  regexp.MustCompilePOSIX("^" + re),
			Handler: h,
		})
	}
	sort.Sort(sort.Reverse(rs))

	nrs := make(PathSwitch, 0, len(routes))
	for _, rr := range rs {
		nrs = append(nrs, PathRule(rr))
	}
	return nrs
}
