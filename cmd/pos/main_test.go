package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectToHTTPSKeepsTheHostAndThePath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://pos.example.org:80/kitchen?a=1", nil)
	recorder := httptest.NewRecorder()

	redirectToHTTPS(recorder, request)

	if recorder.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMovedPermanently)
	}
	got := recorder.Header().Get("Location")
	want := "https://pos.example.org/kitchen?a=1"
	if got != want {
		t.Fatalf("Location = %q, want %q", got, want)
	}
}
