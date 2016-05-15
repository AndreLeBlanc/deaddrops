package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"deadrop/api"
)

func endUpload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	decoder := json.NewDecoder(r.Body)
	meta := decodeJson(decoder)

	fmt.Println(meta)
	fmt.Fprintf(w, "%v", meta.Token)
}

func decodeJson(decoder *json.Decoder) api.Stash {
	var meta api.Stash
	err := decoder.Decode(&meta)
	if err != nil {
		fmt.Printf("the error is ", err)
	}
	return meta
}

func finalize(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "POST" {
		fmt.Println("Finalize: Invalid request")
		return
	}

	t := r.Header.Get("Content-Type")
	if t == "application/json" {
		fmt.Println("I just received a JSON")
	}
	endUpload(w, r, conf)
}
