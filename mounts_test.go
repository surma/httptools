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

func TestMounts_OriginalPath(t *testing.T) {
	h := Mounts{
		"/first/handler": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Path", w.(VarsResponseWriter).Vars()["OrigPath"].(string))
		}),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/first/handler/and/a/path", nil)))
	expected := []string{"/first/handler/and/a/path"}
	expectedCode := http.StatusOK
	got := rr.HeaderMap["X-Path"]
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected error code. Expected %d, got %d", expectedCode, rr.Code)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestMounts_CascadingOriginalPath(t *testing.T) {
	h := Mounts{
		"/first/": Mounts{
			"/handler": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Path", w.(VarsResponseWriter).Vars()["OrigPath"].(string))
			}),
		},
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/first/handler/and/a/path", nil)))
	expected := []string{"/first/handler/and/a/path"}
	expectedCode := http.StatusOK
	got := rr.HeaderMap["X-Path"]
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected error code. Expected %d, got %d", expectedCode, rr.Code)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func ExampleMounts() {
	ms := Mounts{
		"/api/": Mounts{
			"/cars": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Request path:", r.URL.Path)
				fmt.Println("Original path:", w.(VarsResponseWriter).Vars()["OrigPath"].(string))
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
	// Original path: /api/cars/bentley
}
