package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func endUpload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	decoder := json.NewDecoder(r.Body)
	meta := decodeJson(decoder)

	fmt.Println(meta)
	fmt.Fprintf(w, "%v", meta.Token)
}

func decodeJson(decoder *json.Decoder) stash {
	var meta stash
	err := decoder.Decode(&meta)
	if err != nil {
		fmt.Printf("the error is ", err)
	}
	return meta
}

func finalize(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

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
