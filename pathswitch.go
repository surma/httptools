package httptools

import (
	"fmt"
	"net/http"
)

// RegexpRule represents a single rule in a RegexpSwitch.
type PathRule interface {
	// ok is true, if the rule matches.
	// The returned array will be saved to the VarsResponseWriter.
	Match(s string) (submatches []string, ok bool)
	http.Handler
}

// PathSwitch dispatches requests to different handlers depending
// on which rule r.URL.Path matches.
// PathSwitch is a slice of PathRules. They will be checked
// in the order they have been provided. If a rule matches
// (i.e. PathRuke.Matching's ok value is true), the
// Handler will be called and the slice traversal is stopped.
type PathSwitch []PathRule

func (ps PathSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vrw, ok := w.(VarsResponseWriter)
	if !ok {
		vrw = newOurResponseWriter(w)
	}

	for _, rule := range ps {
		if submatches, ok := rule.Match(r.URL.Path); ok {
			for i := 0; i < len(submatches); i++ {
				vrw.Vars()[fmt.Sprintf("%d", i+1)] = submatches[i]
			}
			rule.ServeHTTP(vrw, r)
			return
		}
	}
	http.Error(vrw, "Not found", http.StatusNotFound)
}
