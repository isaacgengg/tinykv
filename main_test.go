package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer(t *testing.T) (*Store, *http.ServeMux) {
	store := NewStore(t.TempDir() + "/wal.log")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /kv/{key}", store.getValue)
	mux.HandleFunc("PUT /kv/{key}", store.putKeyValue)
	mux.HandleFunc("DELETE /kv/{key}", store.deleteKey)
	return store, mux
}

func TestGetValue(t *testing.T) {
	_, mux := newTestServer(t)
	req := httptest.NewRequest("GET", "/kv/foo", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

}

func TestPutAndGetValue(t *testing.T) {
	_, mux := newTestServer(t)

	req := httptest.NewRequest("PUT", "/kv/foo", strings.NewReader("bar"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/kv/foo", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Body.String() != "bar" {
		t.Errorf("expected bar, got %s", w.Body.String())
	}

}

func TestDeleteKey(t *testing.T) {
	_, mux := newTestServer(t)

	req := httptest.NewRequest("PUT", "/kv/foo", strings.NewReader("bar"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	req = httptest.NewRequest("DELETE", "/kv/foo", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	req = httptest.NewRequest("GET", "/kv/foo", nil)
	w = httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
