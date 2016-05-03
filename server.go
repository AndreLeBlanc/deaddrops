package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//global map of channels
var mapOfChannels map[string]chan string

func createStash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		//Create token for upload session
		token := md5.New()
		t := time.Now()
		io.WriteString(token, t.String())
		c := make(chan string)
		//TODO probably will not work with a global variable, use supersupervisor??
		stringToken := hex.EncodeToString(token.Sum(nil))
		appendChan(mapOfChannels, stringToken, c)
		//go supervisor(token, c)
		reply, _ := json.Marshal(stringToken)
		//TODO: handle error from JSON
		fmt.Fprintf(w, string(reply))
	} else if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		//TODO check that token is valid
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./filetest/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	} else {
		return
	}
}

func initServer() {
	//TODO: check/start database

	//TODO: load server settings from somewhere, ex. port number
	mapOfChannels = initChanMap()
	http.HandleFunc("/", createStash)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initChanMap() map[string]chan string {
	//TODO: chan string should maybe be a struct instead with relevant data
	return map[string]chan string{}
}

func appendChan(m map[string]chan string, token string, c chan string) {
	m[token] = c
}

func getChan(m map[string]chan string, token string) chan string {
	return m[token]
}

func main() {
	initServer()
}
