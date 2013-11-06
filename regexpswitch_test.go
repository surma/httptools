package httptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegexpSwitch_IsModifiedResponseWriter(t *testing.T) {
	rs := NewRegexpSwitch(map[string]http.Handler{
		"^.*$": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(VarsResponseWriter)
			w.Header().Add("WasVRW", fmt.Sprintf("%v", ok))
			_, ok = w.(CheckResponseWriter)
			w.Header().Add("WasCRW", fmt.Sprintf("%v", ok))
		}),
	})

	rr := httptest.NewRecorder()
	rs.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/", nil)))
	expected := "true"
	got := rr.HeaderMap.Get("WasVRW")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
	got = rr.HeaderMap.Get("WasCRW")
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestRegexpSwitch_Routing(t *testing.T) {
	rs := NewRegexpSwitch(map[string]http.Handler{
		"^/([a-z]+)/?$": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vrw := w.(VarsResponseWriter)
			w.Header().Set("X-Path", vrw.Vars()["1"].(string))
		}),
	})

	rr := httptest.NewRecorder()
	rs.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/testpath/", nil)))
	expected := []string{"testpath"}
	got := rr.HeaderMap["X-Path"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	rs.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/!!!/", nil)))
	expectedCode := http.StatusNotFound
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected status code. Expected %d, got %d", expectedCode, rr.Code)
	}
}

func ExampleNewRegexpSwitch() {
	rr := NewRegexpSwitch(map[string]http.Handler{
		"/people/([a-z]+)": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := w.(VarsResponseWriter).Vars()
			fmt.Printf("You are looking for %s", vars["1"].(string))
		}),
	})
	req, _ := http.NewRequest("GET", "/people/peter", nil)
	rr.ServeHTTP(httptest.NewRecorder(), req)
	// Output:
	// You are looking for peter
}
