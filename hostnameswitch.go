package httptools

import (
	"net/http"
	"regexp"
)

// HostnameSwitch dispatches requests to different handlers depending
// on the value of r.Host. If no appropriate handler is found, the handler
// with the "_" key will be used. Otherwise, a 404 is returned.
// Port numbers in the request will be stripped before matching.
type HostnameSwitch map[string]http.Handler

var (
	portStripper = regexp.MustCompilePOSIX(":[0-9]+$")
)

func (hs HostnameSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := portStripper.ReplaceAllString(r.Host, "")
	h, ok := hs[host]
	if ok {
		h.ServeHTTP(w, r)
		return
	}
	h, ok = hs["_"]
	if ok {
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}
