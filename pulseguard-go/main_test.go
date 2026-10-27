package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestValidURL(t *testing.T) {
	for _, tc := range []struct {
		u    string
		want bool
	}{{"https://example.org", true}, {"http://localhost:8080", true}, {"file:///etc/passwd", false}, {"https://user:pass@host", false}, {"not-a-url", false}} {
		if got := validURL(tc.u); got != tc.want {
			t.Errorf("%s got %v", tc.u, got)
		}
	}
}
func TestCheckSite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	defer srv.Close()
	got := checkSite(context.Background(), srv.Client(), srv.URL)
	if !got.Online || got.Status != 204 {
		t.Fatalf("unexpected check: %+v", got)
	}
}
func TestPersistence(t *testing.T) {
	f := filepath.Join(t.TempDir(), "sites.json")
	s, _ := loadStore(f)
	a, err := s.add("Demo", "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	s2, err := loadStore(f)
	if err != nil || len(s2.all()) != 1 {
		t.Fatalf("failed to persist: %v", err)
	}
	if !s2.remove(a.ID) {
		t.Fatal("could not remove")
	}
	if len(s2.all()) != 0 {
		t.Fatal("still present")
	}
}
