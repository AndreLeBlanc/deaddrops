package server

import (
	"net/http"
	//"net/http/httptest"
	"testing"
)

func TestSysSuper(t *testing.T) {

}

func TestUpSuperUpload(t *testing.T) {
	conf := InitServer()
	w, _ := createGet("http://localhost:9090/upload", conf)
	if w == nil {
		t.Errorf("Error creating GET [Create] request")
	}
	if w.Code != http.StatusOK {
		t.Errorf("[Create] GET response error: %v", w.Code)
	}

	token, err := getToken(w)
	if err != nil {
		t.Errorf("[Upload] Invalid token format")
	}

	reply, err := UpSuperUpload(token, "test1.txt", conf)
	if err != nil {
		t.Errorf("[UpSuper] Invalid token")
	}
	if len(reply.Meta.Files) != 1 {
		t.Errorf("[UpSuper] Did not append to Files array")
	}

	if reply.Meta.Files[0].Fname != "test1.txt" {
		t.Errorf("[UpSuper] Wrong file in Files array")
	}
}

func TestUpSuperFinalize(t *testing.T) {

}

func TestDnSuperStash(t *testing.T) {

}

func TestDnSuperDownload(t *testing.T) {

}
