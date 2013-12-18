package httptools

import (
	"net/http"
)

// VarsResponseWriter is a http.ResponseWriter which gives access
// to a map. The map can be filled with arbitrary data and is supposed
// to be out-of-band channel to pass data between handlers in a handler list
// or any kind of handler switch.
type VarsResponseWriter interface {
	http.ResponseWriter
	Vars() map[string]interface{}
}

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
	if vrw, ok := w.(VarsResponseWriter); ok {
		orw.vars = vrw.Vars()
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
