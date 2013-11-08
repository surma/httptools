package httptools

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHostnameSwitch_NoDefault(t *testing.T) {
	h := HostnameSwitch{
		"www.google.com": http.HandlerFunc(handlerA),
		"www.google.de":  http.HandlerFunc(handlerB),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "http://www.google.com", nil)))
	expected := []string{"a"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "http://www.google.de", nil)))
	expected = []string{"b"}
	got = rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "http://www.google.it", nil)))
	expectedCode := http.StatusNotFound
	if rr.Code != expectedCode {
		t.Fatalf("Unexpected error code. Expected %d, got %d", expectedCode, rr.Code)
	}
}

func TestHostnameSwitch_Default(t *testing.T) {
	h := HostnameSwitch{
		"www.google.com": http.HandlerFunc(handlerA),
		"_":              http.HandlerFunc(handlerB),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "http://www.google.com", nil)))
	expected := []string{"a"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "http://www.google.de", nil)))
	expected = []string{"b"}
	got = rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}

func TestHostnameSwitch_IgnorePortNumbers(t *testing.T) {
	h := HostnameSwitch{
		"www.google.com": http.HandlerFunc(handlerA),
		"_":              http.HandlerFunc(handlerB),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, MustRequest(http.NewRequest("GET", "http://www.google.com:8080/Test/123", nil)))
	expected := []string{"a"}
	got := rr.HeaderMap["Handler"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Header list wrong. Expected %#v, got %#v", expected, got)
	}
}
