package main

import (
	"fmt"
	"net/http"
	"io"
	"log"
	"os"
)

func createStash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		//http.ServeFile(w,r, "./filetest/README.md")
	} else if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
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

func appendChan(m map[string]chan string, token string) {
	m[token] = make(chan string)
}

func getChan(m map[string]chan string, token string) chan string {
	return m[token]
}

func main() {
	initServer()
}
