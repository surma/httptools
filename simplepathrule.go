package httptools

import (
	"net/http"
	"regexp"
	"sort"
)

// SimplePathRule is a PathRule using a common placeholder syntax.
// The prefix is matched verbatim, names in braces catch any symbol
// except `/` and are returned as a submatch
type SimplePathRule struct {
	Rule string
	http.Handler
	re *regexp.Regexp
}

var converterRegexp = regexp.MustCompilePOSIX("\\{[^\\}]+\\}")

func (spr SimplePathRule) Match(s string) ([]string, bool) {
	if spr.re == nil {
		re := converterRegexp.ReplaceAllString(spr.Rule, "([^/]+)")
		spr.re = regexp.MustCompilePOSIX("^" + re)
	}
	submatches := spr.re.FindStringSubmatch(s)
	if submatches == nil || len(submatches) < 1 {
		return nil, false
	}
	return submatches[1:], true
}

func (spr SimplePathRule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if spr.Handler != nil {
		spr.Handler.ServeHTTP(w, r)
	}
}

type simpleSwitch []SimplePathRule

func (ss simpleSwitch) Len() int {
	return len(ss)
}

func (ss simpleSwitch) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func (ss simpleSwitch) Less(i, j int) bool {
	return len(ss[i].re.String()) < len(ss[j].re.String())
}

// A regexp switch takes a map of regexp strings and handlers.
// A regexp is considered a match if it matches the string from the
// beginning (but not necessarily to the end).
// Longer patterns take precedence over shorter ones.
func NewSimpleSwitch(routes map[string]http.Handler) PathSwitch {
	ss := make(simpleSwitch, 0, len(routes))
	for rule, h := range routes {
		ss = append(ss, SimplePathRule{
			Rule:    rule,
			Handler: h,
		})
	}
	sort.Sort(sort.Reverse(ss))

	ps := make(PathSwitch, 0, len(routes))
	for _, sr := range ss {
		ps = append(ps, PathRule(sr))
	}
	return ps
}
