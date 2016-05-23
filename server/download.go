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
		http.Error(w, "Invalid request", 400)
		return
	}

	urlSubStr := api.ParseURL(r.URL.Path)
	if len(urlSubStr) < 3 {
		http.Error(w, "Invalid URL", 400)
		return
	}
	if len(urlSubStr) == 3 {
		if api.ValidateToken(api.GetToken(urlSubStr)) {
			createJsonStash(w, api.GetToken(urlSubStr), conf)
			return
		} else {
			http.Error(w, "Invalid URL", 400)
			return
		}
	}
	if !api.ValidateFileId(api.GetFilename(urlSubStr)) || !api.ValidateToken(api.GetToken(urlSubStr)) {
		http.Error(w, "Invalid Filename", 404)
		return
	}

	token := api.GetToken(urlSubStr)
	fileid := api.GetFilename(urlSubStr)

	reply, err := DnSuperDownload(token, fileid, conf)
	if err != nil || reply.HttpCode != http.StatusOK {
		http.Error(w, reply.Message, reply.HttpCode)
		return
	}
	filename := reply.Meta.Files[0].Fname
	path := filepath.Join(conf.filefolder, token, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println(err)
		http.Error(w, "File does not exist", 404)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename='" + filename + "'")
	http.ServeFile(w, r, path)
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
		http.Error(w, "Internal server error", 500)
		return
	}

	w.Write([]byte(json))
}
