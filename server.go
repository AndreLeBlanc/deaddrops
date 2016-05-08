package main

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
	"regexp"
	"time"
	"path/filepath"
)

//TODO temporary config fake struct
type Configuration struct {
	filefolder string
	port       string
}

func createStash(w http.ResponseWriter, r *http.Request, cm *api.ChanMap, conf *Configuration) {
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
		api.AppendChan(cm, stringToken, c)

		//go supervisor(token, c, cm) //TODO: maybe skip c?
		go dummySupervisor(stringToken, c, cm)

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

		validateFile( /*file*/)
		
		if _, err := os.Stat(filepath.Join(conf.filefolder,token)); os.IsNotExist(err){
			err2:= os.MkdirAll(filepath.Join(conf.filefolder, token), 0700 )
			fmt.Println("Creating new token folder")
			if err2 != nil{
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
		fmt.Println("wrote to file")
		//TODO: maybe have a response channel for the supervisor to reply
		//ie. c <- handler.Filename, responseChannel
		c <- handler.Filename
		fmt.Println("Sent filename to channel")
		fmt.Fprintf(w, "%v", handler.Header)
	} else {
		return
	}
}

func download(w http.ResponseWriter, r *http.Request, cm *api.ChanMap, conf *Configuration) {
	fmt.Println("method:", r.Method)
	if r.Method != "GET" {
		// Invalid request
		return
	}

	token := r.FormValue("token")
	filename := r.FormValue("filename")
	filepath := filepath.Join(conf.filefolder, token, filename)

	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Disposition", "attachment; filename='"+filename+"'")
	http.ServeFile(w, r, filepath)
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

func dummySupervisor(token string, c chan string, cm *api.ChanMap) {
	fmt.Println("Upload supervisor %s up and running", token)
	loop := true
	for loop {
		select {
		case fname := <-c:
			fmt.Println("received filename: %s", fname)
			//	case <-time.After(time.Second + 100)://TODO decide timeout
			//		fmt.Println("timeout")
		}
	}
}

var validPath = regexp.MustCompile("^/(test|download)")

func makeHandler(f func(http.ResponseWriter, *http.Request, *api.ChanMap, *Configuration), cm *api.ChanMap, conf *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			fmt.Println("invalid path")
			return
		}
		f(w, r, cm, conf)
	}
}

func (c *Configuration) loadSettings() {
	//TODO: load server settings from somewhere, ex. port number
	c.filefolder = "deadropfiles"
	c.port = ":8080"
}

func initServer() {
	//TODO: check/start database
	var conf = new(Configuration)
	conf.loadSettings()
	cm := api.InitChanMap()
	//Check if folder "deadropfiles" exist
	if _, err := os.Stat(conf.filefolder); os.IsNotExist(err) {
		err2 := os.Mkdir(conf.filefolder, 0700) //Borde det vara 0700?
		fmt.Println("Creating folder ", conf.filefolder)
		if err2 != nil {
			log.Fatal("Could not create file directory %s", err2)
		}
	} else {
		fmt.Println("Folder exists")
	}
	http.HandleFunc("/test", makeHandler(createStash, cm, conf))
	http.HandleFunc("/download", makeHandler(download, cm, conf))

	err := http.ListenAndServe(conf.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func validateFile( /*file*/) bool {
	//TODO: file validation, ex. not too big
	return true
}

func main() {
	initServer()
}
