package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlerList_Order(t *testing.T) {
	h := HandlerList{
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
	h := HandlerList{
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
