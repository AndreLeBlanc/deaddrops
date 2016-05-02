package main

import (
	"fmt"
	"net/http"
//	"html/template"
//	"strconv"
	"log"
	"io"
	"os"
//	"crypto/md5"
//	"time"
)

func uploadFile(w http.ResponseWriter, r *http.Request){
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
	    http.ServeFile(w,r, "./filetest/README.md")
    } else {
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
    }
}

func main(){
	http.HandleFunc("/upload", uploadFile)
	err:= http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
