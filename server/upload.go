package server

import (
	"crypto/md5"
	"deadrop/api"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

func createStash(w http.ResponseWriter, r *http.Request, cm *api.ChanMap) {
	//Create token for upload session
	token := md5.New()
	t := time.Now()
	io.WriteString(token, t.String())
	//TODO probably will not work with a global variable, use supersupervisor??
	stringToken := hex.EncodeToString(token.Sum(nil))
	c := make(chan string)
	api.AppendChan(cm, stringToken, c)

	//go supervisor(token, c, cm) //TODO: maybe skip c?
	go api.DummySupervisor(stringToken, c, cm)

	reply, _ := json.Marshal(stringToken)
	//TODO: handle error from JSON
	fmt.Fprintf(w, string(reply))
}

func uploadFile(w http.ResponseWriter, r *http.Request, cm *api.ChanMap) {
	//TODO: handle json form at the end, ie. they will send a json object instead of a file

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	token := r.FormValue("token")

	if !validateToken(token) {
		//Abandon ship
		return
	}
	fmt.Println("checked that token is valid")

	c, ok := api.FindChan(cm, token)
	if !ok {
		return
	}
	fmt.Println("Checked that channel exist")

	validateFile( /*file*/ )

	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./filetest/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Println("wrote to file")
	//TODO: maybe have a response channel for the supervisor to reply
	//ie. c <- handler.Filename, responseChannel
	c <- handler.Filename
	fmt.Println("Sent filename to channel")
	fmt.Fprintf(w, "%v", handler.Header)
}

func upload(w http.ResponseWriter, r *http.Request, cm *api.ChanMap) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		createStash(w, r, cm)
	} else if r.Method == "POST" {
		uploadFile(w, r, cm)
	} else {
		return
	}
}

func validateToken(token string) bool {
	if len(token) != 32 {
		return false
	} else if match, _ := regexp.MatchString("^[a-zA-Z0-9]*$", token); !match {
		return match
	} else {
		return true
	}

}

func validateFile( /*file*/ ) bool {
	//TODO: file validation, ex. not too big
	return true
}
