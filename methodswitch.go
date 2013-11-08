package httptools

import (
	"net/http"
)

// MethodSwitch dispatches request to different handlers depending
// on the HTTP verb used in the request.
type MethodSwitch map[string]http.Handler

func (ms MethodSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := ms[r.Method]
	if !ok {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if handler != nil {
		handler.ServeHTTP(w, r)
	}
}
