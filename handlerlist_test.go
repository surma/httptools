package handlerlist

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlerList_IsVarsResponseWriter(t *testing.T) {
	h := L{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(VarsResponseWriter)
			w.Header().Add("WasVRW", fmt.Sprintf("%v", ok))
		}),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := "true"
	got := rr.HeaderMap.Get("WasVRW")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHandlerList_IsSilentVarsResponseWriter(t *testing.T) {
	h := L{
		SilentHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(VarsResponseWriter)
			w.Header().Add("WasVRW", fmt.Sprintf("%v", ok))
		})),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := "true"
	got := rr.HeaderMap.Get("WasVRW")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHandlerList_VarsResponseWriterPersistency(t *testing.T) {
	h := L{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.(VarsResponseWriter).Vars()["SomeData"] = "Data"
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok := w.(VarsResponseWriter).Vars()["SomeData"].(string) == "Data"
			w.Header().Add("WasVRWDataCorrect", fmt.Sprintf("%v", ok))
		}),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	expected := "true"
	got := rr.HeaderMap.Get("WasVRWDataCorrect")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHandlerList_Order(t *testing.T) {
	h := L{
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
	h := L{
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

func TestMethodSwitch(t *testing.T) {
	h := MethodSwitch{
		"GET":  http.HandlerFunc(handlerA),
		"POST": http.HandlerFunc(handlerB),
		"PUT":  http.HandlerFunc(handlerC),
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/", nil)))
	expected := []string{"a"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("POST", "/", nil)))
	expected = []string{"b"}
	got = rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("PUT", "/", nil)))
	expected = []string{"c"}
	got = rr.HeaderMap["Handler"]
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
