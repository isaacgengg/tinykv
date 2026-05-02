package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetValue(t *testing.T) {
	t.Cleanup(func() {
		kvStore = make(map[string]string)
	})

	req := httptest.NewRequest("GET", "/kv/foo", nil)
	w := httptest.NewRecorder()

	getValue(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

}

func TestPutAndGetValue(t *testing.T) {
	t.Cleanup(func() {
		kvStore = make(map[string]string)
	})

	req := httptest.NewRequest("PUT", "/kv/foo", strings.NewReader("bar"))
	w := httptest.NewRecorder()

	putKeyValue(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/kv/foo", nil)
	w = httptest.NewRecorder()

	getValue(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	if w.Body.String() != "bar" {
		t.Errorf("expected bar, got %s", w.Body.String())
	}

}

func TestDeleteKey(t *testing.T) {
	t.Cleanup(func() {
		kvStore = make(map[string]string)
	})

	req := httptest.NewRequest("PUT", "/kv/foo", strings.NewReader("bar"))
	w := httptest.NewRecorder()

	putKeyValue(w, req)

	req = httptest.NewRequest("DELETE", "/kv/foo", nil)
	w = httptest.NewRecorder()

	deleteKey(w, req)

	req = httptest.NewRequest("GET", "/kv/foo", nil)
	w = httptest.NewRecorder()

	getValue(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
