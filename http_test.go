package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)


func TestStashGET(t *testing.T) {
	cm := initChanMap()
	csHandler := makeHandler(createStash, cm)
	req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("No Stash GET response: %v", http.StatusOK)
	}
}
