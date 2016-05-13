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
	"time"
	"deadrop/database"
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


	// Ny kod från André

	db := database.NewConnect()


	//

	go database.SupervisorUp(db, /*token*/"pppp", c) //TODO: maybe skip c? tog bort cm
	go database.SupervisorDo(db, stringToken, c) // tog bort cm
	jsonToken := struct {
		Token string
	}{
		stringToken,
	}
	reply, _ := json.Marshal(jsonToken)
	//TODO: handle error from JSON
	w.Header().Set("Content-Type","application/json")
	w.Write(reply)
}

type stashpayload struct {
	Token string
	//endTime time.Time
	//files map[string]int
}

func endUpload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	decoder := json.NewDecoder(r.Body)
	var meta stashpayload
	err := decoder.Decode(&meta)
	if err != nil {
		fmt.Printf("the error is ", err)
	}
	fmt.Printf("the payload is %s", meta.Token)
	fmt.Fprintf(w, "%v", r.Header)
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

	if !api.ValidateToken(token) {
		//Abandon ship
		return
	}
	fmt.Println("Checked that token is valid")

	c, ok := api.FindChan(conf.ChanMap(), token)
	if !ok {
		return
	}
	fmt.Println("Checked that channel exist")

	api.ValidateFile( /*file*/ ) // ?

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
		t := r.Header.Get("Content-Type")
		if t == "application/json" {
			endUpload(w, r, conf)
			fmt.Println("I just received a JSON")
			return
		}
		uploadFile(w, r, conf)
	} else {
		//TODO return error
		return
	}
}
