package httptools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleDiscardPathElements() {
	ms := L{
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
