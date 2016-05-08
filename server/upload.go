package server

import (
	"crypto/md5"
	"deadrop/api"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func createStash(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	//Create token for upload session
	token := md5.New()
	t := time.Now()
	io.WriteString(token, t.String())
	//TODO probably will not work with a global variable, use supersupervisor??
	cm := conf.ChanMap()
	stringToken := hex.EncodeToString(token.Sum(nil))
	c := make(chan string)
	api.AppendChan(cm, stringToken, c)

	//go supervisor(token, c, cm) //TODO: maybe skip c?
	go api.DummySupervisor2(stringToken, c, cm)

	reply, _ := json.Marshal(stringToken)
	//TODO: handle error from JSON
	fmt.Fprintf(w, string(reply))
}

func uploadFile(w http.ResponseWriter, r *http.Request, conf *Configuration) {
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
	fmt.Println("Checked that token is valid")

	c, ok := api.FindChan(conf.ChanMap(), token)
	if !ok {
		return
	}
	fmt.Println("Checked that channel exist")

	validateFile( /*file*/ ) // ?

	if _, err := os.Stat(filepath.Join(conf.filefolder, token)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(conf.filefolder, token), 0700)
		fmt.Println("Creating new token folder")
		if err != nil {
			log.Fatal("Could not create token folder")
		}
	}

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

func upload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		createStash(w, r, conf)
	} else if r.Method == "POST" {
		uploadFile(w, r, conf)
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
