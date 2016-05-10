package server

import (
	"bytes"
	"deadrop/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStashGET(t *testing.T) {
	conf := InitServer()
	cm := conf.ChanMap()

	if n := api.LenChan(cm); n != 0 {
		t.Errorf("Map containing more/less elements than it should: %d elem", api.LenChan(cm))
	}

	for i := 0; i < 10; i++ {
		csHandler := makeHandler(upload, conf)
		req, _ := http.NewRequest("GET", "http://localhost:8080/upload", nil)
		w := httptest.NewRecorder()
		csHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("[%d] No Stash GET response: %v", i, http.StatusOK)
		}
	}

	if n := api.LenChan(cm); n != 10 {
		t.Errorf("Map containing more/less elements than it should: %d elem", api.LenChan(cm))
	}
}

// Test does not work yet, as POST is not 100% defined yet
func TestStashPOST(t *testing.T) {
	conf := InitServer()

	csHandler := makeHandler(upload, conf)

	req, _ := http.NewRequest("GET", "http://localhost:8080/upload", nil)
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, _ = http.NewRequest("POST", "http://localhost:8080/upload", bytes.NewBuffer(jsonStr))

	w = httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("No Stash POST response: %v", http.StatusOK)
	}
}
