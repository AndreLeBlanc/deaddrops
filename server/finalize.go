package server

import (
	"deadrop/api"
	"encoding/json"
	"fmt"
	"net/http"
)

func endUpload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	decoder := json.NewDecoder(r.Body)
	meta := decodeJson(decoder)
	if meta == nil {
		http.Error(w, "Internal server error", 500)
	}

	// TODO: Uncomment when database is in place.
	// reply, err := UpSuperFinalize(*meta, conf)
	// if err != nil {
	// 	http.Error(w, reply.Message, reply.HttpCode)
	// 	return
	// }

	fmt.Println(*meta)
	fmt.Fprintf(w, "%v", meta.Token)
}

func decodeJson(decoder *json.Decoder) *api.Stash {
	var meta api.Stash
	err := decoder.Decode(&meta)
	if err != nil {
		fmt.Printf("the error is ", err)
		return nil
	}
	return &meta
}

func finalize(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "POST" {
		fmt.Println("Finalize: Invalid request")
		http.Error(w, "Invalid request", 400)
		return
	}

	t := r.Header.Get("Content-Type")
	if t == "application/json" {
		fmt.Println("I just received a JSON")
	}
	endUpload(w, r, conf)
}
