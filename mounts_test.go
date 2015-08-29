package httptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMounts_Stripping(t *testing.T) {
	h := Mounts{
		"/first/handler": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Path", r.URL.Path)
		}),
		"/second/handler/": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Path", r.URL.Path)
		}),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/first/handler/and/a/path", nil)))
	expected := []string{"/and/a/path"}
	expectedCode := http.StatusOK
	got := rr.HeaderMap["X-Path"]
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected error code. Expected %d, got %d", expectedCode, rr.Code)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/second/handler/and/a/path", nil)))
	expected = []string{"/and/a/path"}
	expectedCode = http.StatusOK
	got = rr.HeaderMap["X-Path"]
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected error code. Expected %d, got %d", expectedCode, rr.Code)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/third/handler/and/a/path", nil)))
	expectedCode = http.StatusNotFound
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected error code. Expected %d, got %d", expectedCode, rr.Code)
	}
}

func ExampleMounts() {
	ms := Mounts{
		"/api/": Mounts{
			"/cars": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Request path:", r.URL.Path)
			}),
			"/people": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// ...
			}),
		},
	}
	req, _ := http.NewRequest("GET", "/api/cars/bentley", nil)
	ms.ServeHTTP(httptest.NewRecorder(), req)
	// Output:
	// Request path: /bentley
}
