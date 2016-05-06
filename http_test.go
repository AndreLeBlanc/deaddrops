package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStashGET(t *testing.T) {
	cm := initChanMap()

	if n := len(cm.m); n != 0 {
		t.Errorf("Map containing more/less elements than it should: %d elem", len(cm.m))
	}

	for i := 0; i < 10; i++ {
		csHandler := makeHandler(createStash, cm)
		req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
		w := httptest.NewRecorder()
		csHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("[%d] No Stash GET response: %v", i, http.StatusOK)
		}
	}

	if n := len(cm.m); n != 10 {
		t.Errorf("Map containing more/less elements than it should: %d elem", len(cm.m))
	}
}

// Test does not work yet, as POST is not 100% defined yet
func TestStashPOST(t *testing.T) {
	cm := initChanMap()
	csHandler := makeHandler(createStash, cm)

	req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, _ = http.NewRequest("POST", "http://localhost:8080/test", bytes.NewBuffer(jsonStr))

	w = httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("No Stash POST response: %v", http.StatusOK)
	}
}
