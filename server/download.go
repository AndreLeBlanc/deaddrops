package server

import (
	"deadrop/api"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func download(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "GET" {
		fmt.Println("Download: Invalid request")
		http.Error(w, "Invalid request", 400)
		return
	}

	fmt.Println(r.URL.Path)
	urlSubStr := api.ParseURL(r.URL.Path)
	fmt.Println(urlSubStr)
	if len(urlSubStr) < 3 {
		http.Error(w, "Invalid URL", 400)
		return
	}
	if len(urlSubStr) == 3 {
		if api.ValidateToken(api.GetToken(urlSubStr)) {
			// TODO: Uncomment when database is in place.
			createJsonStash(w, api.GetToken(urlSubStr), conf)
			return
		} else {
			http.Error(w, "Invalid URL", 400)
			return
		}
	}
	if !api.ValidateFileName(api.GetFilename(urlSubStr)) || !api.ValidateToken(api.GetToken(urlSubStr)) {
		http.Error(w, "Invalid Filename", 404)
		return
	}

	token := api.GetToken(urlSubStr)
	filename := api.GetFilename(urlSubStr)
	path := filepath.Join(conf.filefolder, token, filename)

	// TODO: Uncomment when database is in place.
	// TODO: Don-t think you have to check the error here.
	reply, _ := DnSuperDownload(token, filename, conf)
	if reply.HttpCode != http.StatusOK {
		http.Error(w, reply.Message, reply.HttpCode)
		return
	}

	fmt.Printf("filename is %s", filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File does not exist", 404)
		return
	}
	fmt.Println(path)
	// TODO: Fix headers
	// w.Header().Set("Content-Type", "multipart/form-data")
	// w.Header().Set("Content-Disposition", "attachment; filename='"+filename+"'")
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
		fmt.Println("Failed token json encoding")
		http.Error(w, "Internal server error", 500)
		return
	}

	w.Write([]byte(json))
}
