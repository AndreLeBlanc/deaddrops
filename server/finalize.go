package server

import (
	"deadrop/api"
	"encoding/json"
	"log"
	"net/http"
)

func endUpload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	decoder := json.NewDecoder(r.Body)
	meta := decodeJson(decoder)
	if meta == nil {
		http.Error(w, "Internal server error", 500)
	}

	reply, err := UpSuperFinalize(*meta, conf)
	if err != nil || reply.HttpCode != http.StatusOK {
		http.Error(w, reply.Message, reply.HttpCode)
		return
	}
	json, _ := json.Marshal(reply.Meta) //Should probably check this error but what the hell
	
	w.Write([]byte(json))
}

func decodeJson(decoder *json.Decoder) *api.Stash {
	var meta api.Stash
	err := decoder.Decode(&meta)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &meta
}

func finalize(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request", 400)
		return
	}

	endUpload(w, r, conf)
}
