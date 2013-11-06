package httptools

import (
	"net/http"
)

// A handler list is a list of http.Handlers which are
// executed sequentially. If a handler is a SilentHandler and
// it produces output (i.e. calls WriteHeader()), it is assumed
// to be an error message/error code and executing the remaining
// handlers in the list will be skipped.
// The ResponseWriter will be an VarsResponseWriter to make data
// passing between handlers more convenient.
type L []http.Handler

func (l L) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orw, ok := w.(*ourResponseWriter)
	if !ok {
		orw = newOurResponseWriter(w)
	}
	for _, h := range l {
		if _, ok := h.(*silentHandler); ok {
			orw.written = false
			h.ServeHTTP(orw, r)
			if orw.WasWritten() {
				break
			}
		} else {
			h.ServeHTTP(orw, r)
		}
	}
}

type silentHandler struct {
	http.Handler
}

// "Casts" the given handler into a silent handler.
// Silent handlers are expected to produce no output. If they
// do, it is assumend to be an error message/error code.
// In a HandlerList, this execution of the list will be aborted if a
// SilentHandler produces output.
func SilentHandler(h http.Handler) http.Handler {
	return &silentHandler{h}
}
