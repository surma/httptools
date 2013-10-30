// Package handlerlist provides the type `L` with which multiple
// http.Handler can be chained to be executed in sequence.
package handlerlist

import (
	"net/http"
)

// A handler list is a list of http.Handlers which are
// executed sequentially. If a handler is a SilentHandler and
// it produces output (i.e. calls WriteHeader()), it is assumed
// to be an error message/error code and executing the remaining
// handlers in the list will be skipped.
type L []http.Handler

func (l L) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = VarsResponseWriter{
		ResponseWriter: w,
		Vars:           map[string]interface{}{},
	}
	for _, h := range l {
		if _, ok := h.(*silentHandler); ok {
			w = &response{w, false}
			h.ServeHTTP(w, r)
			if w.(*response).written {
				break
			}
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

// VarsResponseWriter is a http.ResponseWriter which gives access
// to a map. The map can be filled with arbitrary data and is supposed
// to be out-of-band channel to pass data between handlers in a HandlerList.
type VarsResponseWriter struct {
	http.ResponseWriter
	Vars map[string]interface{}
}

type silentHandler struct {
	http.Handler
}

// "Casts" the given handler into a silent handler.
// Silent handlers are expected to produce no output. If they
// do, it is assumend to be an error message/error code.
// In a HandlerList, this execution of the list will be aborted if a
// SilentHandler produces output.
func SilentHandler(h http.Handler) *silentHandler {
	return &silentHandler{h}
}

// A wrapper for http.ResponseWriter to record
// if a header has been written
type response struct {
	http.ResponseWriter
	written bool
}

func (r *response) WriteHeader(n int) {
	r.written = true
	r.ResponseWriter.WriteHeader(n)
}
