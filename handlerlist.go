// Package handlerlist provides the type `L` with which multiple
// http.Handler can be chained to be executed sequentially.
//
// Example:
//
//    func userData(w http.ResponseWriter, r *http.Request) {
//    	// Session magic
//    	session, err := openSession(r)
//    	if err != nil {
//    		http.Error(w, "Could not open session", http.StatusBadRequest)
//    	}
//    	w.(*handlerlist.VarsResponseWriter).Vars["UID"] = session.UserId
//    }
//
//    func showProfile(w http.ResponseWriter, r *http.Request) {
//    	uid := w.(*handlerlist.VarsResponseWriter).Vars["UID"].(string)
//
//    	profile := userProfile(uid)
//    	renderProfileTemplate(w, profile)
//    }
//
//    func main() {
//    	// ...
//    	http.Handle("/profile", handlerList.L {
//    		http.HandlerFunc(userData),
//    		handlerlist.SilentHandler(
//    			http.HandlerFunc(showProfile),
//    		)
//    	})
//    	// ...
//    }
package handlerlist

import (
	"net/http"
)

const (
	VERSION = "1.0.0"
)

// A handler list is a list of http.Handlers which are
// executed sequentially. If a handler is a SilentHandler and
// it produces output (i.e. calls WriteHeader()), it is assumed
// to be an error message/error code and executing the remaining
// handlers in the list will be skipped.
// The ResponseWriter will have an VarsResponseWriter as an underlying
// type to make data passing between handlers more convenient.
type L []http.Handler

func (l L) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orw := newOurResponseWriter(w)
	for _, h := range l {
		if _, ok := h.(*silentHandler); ok {
			orw.written = false
			h.ServeHTTP(orw, r)
			if orw.written {
				break
			}
		} else {
			h.ServeHTTP(orw, r)
		}
	}
}

// VarsResponseWriter is a http.ResponseWriter which gives access
// to a map. The map can be filled with arbitrary data and is supposed
// to be out-of-band channel to pass data between handlers in a handler list.
type VarsResponseWriter interface {
	http.ResponseWriter
	Vars() map[string]interface{}
}

func newOurResponseWriter(w http.ResponseWriter) *ourResponseWriter {
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

func (orw *ourResponseWriter) WriteHeader(n int) {
	orw.written = true
	orw.ResponseWriter.WriteHeader(n)
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

// MethodSwitch offers a simple way to apply different handlers depending
// on the HTTP verb used in the request.
type MethodSwitch map[string]http.Handler

func (ms MethodSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := ms[r.Method]
	if !ok {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	handler.ServeHTTP(w, r)
}
