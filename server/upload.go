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
	var st []api.StashFile
	s := api.Stash{Token: token, Lifetime:0, Files:append(st, api.StashFile{Fname:handler.Filename,Size:0,Type:"",Download:0})}


	//TODO: maybe have a response channel for the supervisor to reply
	//ie. c <- handler.Filename, responseChannel
	replyChannel := make(chan api.HttpReplyChan)
	c <- api.SuperChan{s,replyChannel}
	fmt.Println("Sent filename to channel")
	supAns := <- replyChannel
	if supAns.HttpCode == 200 {
		f, err := os.OpenFile(filepath.Join(conf.filefolder, token, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Fprintf(w, "%v", handler.Header)
	
	} else {
		http.Error(w, supAns.Message, supAns.HttpCode)
	}
}
