package server

import (
	"deadrop/api"
	//"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func download(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	fmt.Println("method:", r.Method)
	if r.Method != "GET" {
		// Invalid request
		return
	}

	urlSubStr, valid := parseURL(r.URL.Path)
	if !valid {
		token, valid := validateURLpath(r.URL.Path)
		if !valid {
			http.Error(w, "Invalid URL", 404) // TODO: Check error code
			return
		}

		json := createJsonStash(token, conf)
		w.Write([]byte(json))
		return
	}

	token := getToken(urlSubStr)
	filename := getFilename(urlSubStr)
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
	_, ok := api.FindChan(conf.downMap, token)
	if !ok {
		// superChan := make(chan string)
		// api.AppendChan(conf.downMap, token, c)
		// api.DownSupervisor(superChan, conf)
	}
	return token
}

func validateURLpath(path string) (string, bool) {
	valid, err := regexp.Compile("^/(download)/([a-zA-Z0-9]+)$")
	if err != nil {
		fmt.Println(err)
		return path, false
	}

	urlSubStr := valid.FindStringSubmatch(path)
	fmt.Println(urlSubStr)
	if len(urlSubStr) < 2 || !api.ValidateToken(urlSubStr[2]) { //TODO: If not lazy evaluations, there could be a Panic
		return path, false
	}

	return urlSubStr[2], true
}

func parseURL(path string) ([]string, bool) {
	valid, err := regexp.Compile("^/(download)/([\\w]+)/([\\w]+\\.[a-z]+)$")
	if err != nil {
		fmt.Println(err)
		return []string{}, false
	}

	urlSubStr := valid.FindStringSubmatch(path)
	if len(urlSubStr) == 0 {
		return []string{}, false
	}

	return urlSubStr, true
}

func getToken(urlSubStr []string) string {
	return urlSubStr[2]
}

func getFilename(urlSubStr []string) string {
	return urlSubStr[3]
}
