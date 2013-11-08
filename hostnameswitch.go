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

func TrimPortNumber(host string) string {
	return portStripper.ReplaceAllString(host, "")
}

func (hs HostnameSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := TrimPortNumber(r.Host)
	h, ok := hs[host]
	if ok {
		if h != nil {
			h.ServeHTTP(w, r)
		}
		return
	}
	h, ok = hs["_"]
	if ok {
		if h != nil {
			h.ServeHTTP(w, r)
		}
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}
