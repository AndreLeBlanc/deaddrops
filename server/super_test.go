package server

import (
	"net/http"
	//"net/http/httptest"
	"deadrop/api"
	"testing"
)

func TestSysSuper(t *testing.T) {

}

func TestUpSuperUpload(t *testing.T) {
	conf := InitServer()
	w, _ := createGet("http://localhost:9090/upload", conf)
	if w == nil {
		t.Errorf("Error creating GET [Create] request")
		return
	}
	if w.Code != http.StatusOK {
		t.Errorf("[Create] GET response error: %v", w.Code)
		return
	}

	token, err := getToken(w)
	if err != nil {
		t.Errorf("[Upload] Invalid token format")
		return
	}

	reply, err := UpSuperUpload(token, "test1.txt", conf)
	if err != nil {
		t.Errorf("[UpSuper] Invalid token")
		return
	}
	if len(reply.Meta.Files) != 1 {
		t.Errorf("[UpSuper] Did not append to Files array")
		return
	}

	if reply.Meta.Files[0].Fname != "test1.txt" {
		t.Errorf("[UpSuper] Wrong file in Files array")
	}

	for i := 0; i < 10; i++ {
		reply, err := UpSuperUpload(token, "test1.txt", conf)
		if err != nil {
			t.Errorf("[UpSuper] Invalid token")
			return
		}

		if len(reply.Meta.Files) != i+2 {
			t.Errorf("[UpSuper] Did not append to Files array")
			return
		}
	}

	tstash = &reply.Meta
	tconf = conf
}

var tconf *Configuration = nil
var tstash *api.Stash = nil

func TestUpSuperFinalize(t *testing.T) {
	if tconf == nil || tstash == nil {
		t.Errorf("UpSuperUpload failed")
		return
	}
	reply, err := UpSuperFinalize(*tstash, tconf)
	if err != nil || reply.HttpCode != http.StatusOK {
		t.Errorf("UpSuperFinalize failed")
		tconf = nil
		tstash = nil
		return
	}
}

func TestDnSuperStash(t *testing.T) {
	if tconf == nil || tstash == nil {
		t.Errorf("UpSuperUpload or UpSuperFinalize failed")
	}
	reply, err := UpSuperFinalize(*tstash, tconf)
	if err != nil || reply.HttpCode != http.StatusOK {
		t.Errorf("UpSuperFinalize failed")
		tconf = nil
		tstash = nil
		return
	}
}

func TestDnSuperDownload(t *testing.T) {

}
