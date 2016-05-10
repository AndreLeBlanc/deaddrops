package server

import (
	"deadrop/api"
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

	token, valid := validateURLpath(r.URL.Path)
	if !valid {
		http.Error(w, "Invalid URL", 404) // TODO: Check error code
		return
	}

	//token := r.FormValue("token")

	filename := r.FormValue("filename")

	if len(filename) == 0 {
		_, ok := api.FindChan(conf.downMap, token)
		if !ok {
			// superChan := make(chan string)
			// api.AppendChan(conf.downMap, token, c)
			// api.DownSupervisor(superChan, conf)

		}

		return
	}

	if !validateFilename(filename) {
		http.Error(w, "Invalid file name", 404) // TODO: Check error code
	}

	filepath := filepath.Join(conf.filefolder, token, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File does not exist", 404) // TODO: Check error code
	}

	// TODO: Fix headers
	// w.Header().Set("Content-Type", "multipart/form-data")
	// w.Header().Set("Content-Disposition", "attachment; filename='"+filename+"'")
	http.ServeFile(w, r, filepath)
}

func validateURLpath(path string) (string, bool) {
	valid, err := regexp.Compile("^/(download)/([a-zA-Z0-9]+)$")
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	urlSubStr := valid.FindStringSubmatch(path)
	fmt.Println(urlSubStr)
	if len(urlSubStr) < 2 || !api.ValidateToken(urlSubStr[2]) { //TODO: If not lazy evaluations, there could be a Panic
		return "", false
	}

	return urlSubStr[2], true
}

func validateFilename(fname string) bool {
	valid, err := regexp.Compile("^[\\w]+\\.[a-z]+$")
	if err != nil {
		fmt.Println(err)
		return false
	}

	return valid.MatchString(fname)
}
