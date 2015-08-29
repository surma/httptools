package httptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDiscardPathElements_TrailingSlash(t *testing.T) {
	ms := List{
		DiscardPathElements(2),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Path", r.URL.Path)
		}),
	}

	rr := httptest.NewRecorder()
	ms.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/some/thing/different/yet/", nil)))
	expected := []string{"/different/yet/"}
	got := rr.HeaderMap["X-Path"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func ExampleDiscardPathElements() {
	ms := List{
		DiscardPathElements(2),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL.Path)
		}),
	}
	req, _ := http.NewRequest("GET", "/prefix/and/a/real/path", nil)
	ms.ServeHTTP(httptest.NewRecorder(), req)
	req, _ = http.NewRequest("GET", "/", nil)
	ms.ServeHTTP(httptest.NewRecorder(), req)
	// Output:
	// /a/real/path
	// /
}
