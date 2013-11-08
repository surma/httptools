package httptools

import (
	"net/http"
)

// HostnameSwitch dispatches requests to different handlers depending
// on the value of r.Host. If no appropriate handler is found, the handler
// with the "_" key will be used. Otherwise, a 404 is returned.
type HostnameSwitch map[string]http.Handler

func (hs HostnameSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := hs[r.Host]
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
