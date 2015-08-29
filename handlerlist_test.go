package httptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlerList_IsModifiedResponseWriter(t *testing.T) {
	h := List{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(CheckResponseWriter)
			w.Header().Add("WasCRW", fmt.Sprintf("%v", ok))
		}),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := "true"
	got := rr.HeaderMap.Get("WasCRW")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHandlerList_IsSilentModifiedResponseWriter(t *testing.T) {
	h := List{
		SilentHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(CheckResponseWriter)
			w.Header().Add("WasCRW", fmt.Sprintf("%v", ok))
		})),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := "true"
	got := rr.HeaderMap.Get("WasCRW")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHandlerList_Order(t *testing.T) {
	h := List{
		SilentHandler(http.HandlerFunc(handlerA)),
		SilentHandler(http.HandlerFunc(handlerB)),
		SilentHandler(http.HandlerFunc(handlerC)),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := []string{"a", "b", "c"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHandlerList_Fail(t *testing.T) {
	h := List{
		SilentHandler(http.HandlerFunc(handlerA)),
		SilentHandler(http.HandlerFunc(handlerB)),
		SilentHandler(http.HandlerFunc(failHandler)),
		SilentHandler(http.HandlerFunc(handlerC)),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := []string{"a", "b"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func handlerA(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler", "a")
}

func handlerB(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler", "b")
}

func handlerC(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler", "c")
}

func failHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Error", http.StatusInternalServerError)
}

func MustRequest(r *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}
	return r
}
