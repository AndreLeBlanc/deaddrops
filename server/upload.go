package server

import (
	"deadrop/api"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func upload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "POST" {
		fmt.Println("Upload: Invalid request")
		http.Error(w, "Invalid request", 400)
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Received bad file", 400)
		return
	}
	defer file.Close()

	token := r.FormValue("token")
	fmt.Printf("[Upload] token: %s\n", token)
	
	if !api.ValidateToken(token) {
		fmt.Println("Invalid token")
		http.Error(w, "Invalid token, bad format", 400)
		return
	}
	fmt.Println("Checked that token is valid")

	c, ok := api.FindChan(conf.upMap, token)
	if !ok {
		fmt.Println("Invalid token, could not find in upMap")
		http.Error(w, "Invalid token, does not exist", 400)
		return
	}
	fmt.Println("Checked that channel exist")

	// TODO: Everything till next comment should be in ValidateFile
	api.ValidateFile( /*file*/ ) // ?
	if _, err := os.Stat(filepath.Join(conf.filefolder, token)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(conf.filefolder, token), 0700)
		fmt.Println("Creating new token folder")
		if err != nil {
			log.Fatal("Could not create token folder")
			http.Error(w, "Internal server error", 500)
			return
		}
	}
	// TODO: end
	filename := ParseFilename(handler.Filename)

	f, err := os.OpenFile(filepath.Join(conf.filefolder, token, filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err) // Could not open file
		http.Error(w, "Internal server error", 500)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//TODO: maybe have a response channel for the supervisor to reply
	//ie. c <- handler.Filename, responseChannel
	c <- filename
	fmt.Println("Sent filename to channel")
	fmt.Fprintf(w, "%v", handler.Header)
}


func ParseFilename(path string) string {
	substr := api.ParseURL(path)
	if len(substr) == 0 {
		fmt.Println("Failed parsing filename")
		return path
	}
	return substr[len(substr)-1] 
}
