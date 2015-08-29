package httptools

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegexpSwitch_Routing(t *testing.T) {
	rs := NewRegexpSwitch(map[string]http.Handler{
		"/([a-z]+)/?": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Flag", "SET")
		}),
	})

	rr := httptest.NewRecorder()
	rs.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/testpath/", nil)))
	expected := []string{"SET"}
	got := rr.HeaderMap["X-Flag"]
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

func TestRegexpSwitch_PatternPrecedence(t *testing.T) {
	rs := NewRegexpSwitch(map[string]http.Handler{
		"/.+":      http.HandlerFunc(handlerA),
		"/some/.+": http.HandlerFunc(handlerB),
	})

	rr := httptest.NewRecorder()
	rs.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/some/thing/bla", nil)))
	expected := []string{"b"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestRegexpSwitch_TrailingSlash(t *testing.T) {
	rs := NewRegexpSwitch(map[string]http.Handler{
		"/(.+)": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Flag", "SET")
		}),
	})

	rr := httptest.NewRecorder()
	rs.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "/some/thing/", nil)))
	expected := []string{"SET"}
	got := rr.HeaderMap["X-Flag"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}
