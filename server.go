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
	"sync"
)


func createStash(w http.ResponseWriter, r *http.Request, cm *ChanMap) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		//Create token for upload session
		token := md5.New()
		t := time.Now()
		io.WriteString(token, t.String())
		//TODO probably will not work with a global variable, use supersupervisor??
		stringToken := hex.EncodeToString(token.Sum(nil))
		c := make(chan string)
		appendChan(cm, stringToken, c)
		//go supervisor(token, c, cm) //TODO: maybe skip c?
		reply, _ := json.Marshal(stringToken)
		//TODO: handle error from JSON
		fmt.Fprintf(w, string(reply))
	} else if r.Method == "POST" {
		//TODO: handle json form at the end, ie. they will send a json object instead of a file
		
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		//TODO: check that token is valid
		token := r.FormValue("token")
		c, ok := findChan(cm, token)
		if !ok {
			//ABANDON SHIP
			return
		}

		validateFile(/*file*/)
		
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./filetest/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		//TODO: maybe have a response channel for the supervisor to reply
		//ie. c <- handler.Filename, responseChannel
		c <- handler.Filename
	} else {
		return
	}
}

func initServer() {
	//TODO: check/start database

	//TODO: load server settings from somewhere, ex. port number
	chanMap := initChanMap()
	//http.HandleFunc("/", createStash)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		createStash(w, r, chanMap)
       })

	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


func validateFile(/*file*/) bool {
	//TODO: file validation, ex. not too big
	return true
}


type ChanMap struct {
	m   map[string]chan string
	mux sync.Mutex
}


func initChanMap() *ChanMap {
	return &ChanMap{m: make(map[string]chan string)}
}

func appendChan(cm *ChanMap, token string, c chan string) {
	cm.mux.Lock()
	if _, ok := findChan(cm, token); ok {
		cm.mux.Unlock()
		return
	}
	cm.m[token] = c
	cm.mux.Unlock()
}

func getChan(cm *ChanMap, token string) chan string {
	return cm.m[token]
}


func findChan(cm *ChanMap, token string) (chan string, bool) {
	c, ok := cm.m[token]
	return c, ok
}


func main() {
	initServer()
}
