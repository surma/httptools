package httptools

import (
	"net/http"
)

// VarsResponseWriter is a http.ResponseWriter which gives access
// to a map. The map can be filled with arbitrary data and is supposed
// to be out-of-band channel to pass data between handlers in a handler list.
type VarsResponseWriter interface {
	http.ResponseWriter
	Vars() map[string]interface{}
}

type CheckResponseWriter interface {
	http.ResponseWriter
	WasWritten() bool
}

func newOurResponseWriter(w http.ResponseWriter) *ourResponseWriter {
	vrw, ok := w.(VarsResponseWriter)
	if ok {
		return &ourResponseWriter{
			ResponseWriter: w,
			vars:           vrw.Vars(),
			written:        false,
		}
	}
	return &ourResponseWriter{
		ResponseWriter: w,
		vars:           map[string]interface{}{},
		written:        false,
	}
}

type ourResponseWriter struct {
	http.ResponseWriter
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
