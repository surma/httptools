package httptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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

func ExampleMethodSwitch() {
	ms := MethodSwitch{
		"GET": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("A GET request")
		}),
		"POST": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("A POST request")
		}),
	}
	req, _ := http.NewRequest("GET", "/", nil)
	ms.ServeHTTP(httptest.NewRecorder(), req)
	req, _ = http.NewRequest("POST", "/", nil)
	ms.ServeHTTP(httptest.NewRecorder(), req)
	// Output:
	// A GET request
	// A POST request
}
