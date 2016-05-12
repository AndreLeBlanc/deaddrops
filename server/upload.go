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
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

	if r.Method != "POST" {
		fmt.Println("Upload: Invalid request")
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	token := r.FormValue("token")
	fmt.Printf("[Upload] token: %s\n", token)
	
	if !api.ValidateToken(token) {
		fmt.Println("Invalid token")
		return
	}
	fmt.Println("Checked that token is valid")

	c, ok := api.FindChan(conf.upMap, token)
	if !ok {
		fmt.Println("Invalid token, could not find in upMap")
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
		}
	}
	// TODO: end

	f, err := os.OpenFile(filepath.Join(conf.filefolder, token, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//TODO: maybe have a response channel for the supervisor to reply
	//ie. c <- handler.Filename, responseChannel
	c <- handler.Filename
	fmt.Println("Sent filename to channel")
	fmt.Fprintf(w, "%v", handler.Header)
}
