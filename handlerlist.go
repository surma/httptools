package httptools

import (
	"net/http"
)

// CheckResponseWriter is a http.ResponseWriter which saves wether it
// has been written to or not.
type CheckResponseWriter interface {
	http.ResponseWriter
	// Returns true if the headers have been written
	WasWritten() bool
}

func newOurResponseWriter(w http.ResponseWriter) *ourResponseWriter {
	orw := &ourResponseWriter{
		ResponseWriter: w,
		vars:           map[string]interface{}{},
		written:        false,
	}
	if hijacker, ok := w.(http.Hijacker); ok {
		orw.Hijacker = hijacker
	}
	return orw
}

type ourResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	vars    map[string]interface{}
	written bool
}

func (orw *ourResponseWriter) Vars() map[string]interface{} {
	return orw.vars
}

func (orw *ourResponseWriter) WasWritten() bool {
	return orw.written
}

func (orw *ourResponseWriter) WriteHeader(n int) {
	orw.written = true
	orw.ResponseWriter.WriteHeader(n)
}

// A handler list is a list of http.Handlers which are
// executed sequentially. If a handler is a SilentHandler and
// it produces output (i.e. calls WriteHeader()), it is assumed
// to be an error message/error code and executing the remaining
// handlers in the list will be skipped.
type List []http.Handler

func (l List) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orw, ok := w.(*ourResponseWriter)
	if !ok {
		orw = newOurResponseWriter(w)
	}
	for _, h := range l {
		if _, ok := h.(*silentHandler); ok {
			orw.written = false
			if h != nil {
				h.ServeHTTP(orw, r)
			}
			if orw.WasWritten() {
				break
			}
		} else {
			if h != nil {
				h.ServeHTTP(orw, r)
			}
		}
	}
}

type silentHandler struct {
	http.Handler
}

// "Casts" the given handler into a silent handler.
// Silent handlers are expected to produce no output. If they
// do, it is assumend to be an error message/error code.
// In a HandlerList, the execution of the list will be aborted if a
// SilentHandler produces output.
func SilentHandler(h http.Handler) http.Handler {
	return &silentHandler{h}
}

// Same as SilentHandler but for http.HandlerFunc
func SilentHandlerFunc(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return &silentHandler{http.HandlerFunc(h)}
}
