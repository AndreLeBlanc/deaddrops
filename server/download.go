package server

import (
	"deadrop/api"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func download(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	urlSubStr := api.ParseURL(r.URL.Path)
	if len(urlSubStr) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	if len(urlSubStr) == 3 {
		if api.ValidateToken(api.GetToken(urlSubStr)) {
			createJsonStash(w, api.GetToken(urlSubStr), conf)
			return
		} else {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
	}
	if !api.ValidateFileId(api.GetFilename(urlSubStr)) || !api.ValidateToken(api.GetToken(urlSubStr)) {
		http.Error(w, "Invalid Filename", http.StatusNotFound)
		return
	}

	token := api.GetToken(urlSubStr)
	fileid := api.GetFilename(urlSubStr)

	reply, err := DnSuperDownload(token, fileid, conf)
	if err != nil || (reply.HttpCode != http.StatusOK && reply.HttpCode != http.StatusResetContent) {
		http.Error(w, reply.Message, reply.HttpCode)
		return
	}

	filearr := reply.Meta.Files
	if len(filearr) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	filename := filearr[0].Fname
	path := filepath.Join(conf.filefolder, token, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println(err)
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
	if reply.HttpCode == http.StatusResetContent {		
		go RmFile(token, filename, conf)
	}
}

func createJsonStash(w http.ResponseWriter, token string, conf *Configuration) {
	reply, err := DnSuperStash(token, conf)
	if err != nil {
		http.Error(w, reply.Message, reply.HttpCode)
		return
	}

	json, err := json.Marshal(reply.Meta)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(json))
}
