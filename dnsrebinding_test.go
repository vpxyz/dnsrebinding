package dnsrebinding

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// a simple handler
var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte("test"))
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
})

func assertResponse(t *testing.T, res *httptest.ResponseRecorder, responseCode int) {
	if responseCode != res.Code {
		t.Errorf("expected response code to be %d but got %d. ", responseCode, res.Code)
	}
}

func TestBadHost(t *testing.T) {
	f := Filter("example.com")

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Host", "malicius.com")

	f(testHandler).ServeHTTP(res, req)

	assertResponse(t, res, 501)
}

func TestGoodHost(t *testing.T) {
	f := Filter("example.com")

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Host", "example.com")

	f(testHandler).ServeHTTP(res, req)

	assertResponse(t, res, 200)
}

func TestSomeGoodHost(t *testing.T) {
	f := Filters("example.com", "google.com", "facebook.com")

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Host", "example.com")

	f(testHandler).ServeHTTP(res, req)

	assertResponse(t, res, 200)
}

func TestSomeBadHost(t *testing.T) {
	f := Filters("example.com", "google.com", "facebook.com")

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Host", "malicius.com")

	f(testHandler).ServeHTTP(res, req)

	assertResponse(t, res, 501)
}
