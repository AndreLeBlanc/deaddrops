package server

import (
	"deadrop/api"
	//"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func download(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "GET" {
		fmt.Println("Download: Invalid request")
		return
	}

	fmt.Println(r.URL.Path)
	urlSubStr := api.ParseURL(r.URL.Path)
	fmt.Println(urlSubStr)
	if len(urlSubStr) < 3 {
		http.Error(w, "Invalid URL", 404) //TODO fix error code
		return
	}
	if len(urlSubStr) == 3 {
		if api.ValidateToken(api.GetToken(urlSubStr)) {
			// Send info about stash to client
			json := createJsonStash(api.GetToken(urlSubStr), conf)
			w.Write([]byte(json))
			return
		} else {
			http.Error(w, "Invalid URL", 404) //TODO fix error code
			return
		}
	}
	if !api.ValidateFileName(api.GetFilename(urlSubStr))||!api.ValidateToken(api.GetToken(urlSubStr)) {
		http.Error(w, "Invalid Filename", 404) //TODO fix error code
		return
	}

	token := api.GetToken(urlSubStr)
	filename := api.GetFilename(urlSubStr)
	path := filepath.Join(conf.filefolder, token, filename)

	fmt.Printf("filename is %s", filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File does not exist", 404) // TODO: Check error code
		return
	}
	fmt.Println(path)
	// TODO: Fix headers
	// w.Header().Set("Content-Type", "multipart/form-data")
	// w.Header().Set("Content-Disposition", "attachment; filename='"+filename+"'")
	http.ServeFile(w, r, path)
}

func createJsonStash(token string, conf *Configuration) string {
	// TODO: This whole function
	_, ok := api.FindChan(conf.downMap, token)
	if !ok {
		// superChan := make(chan string)
		// api.AppendChan(conf.downMap, token, c)
		// api.DownSupervisor(superChan, conf)
	}
	return token
}
