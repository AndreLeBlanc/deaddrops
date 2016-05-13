package server

import (
	"deadrop/api"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

//TODO temporary config fake struct
type Configuration struct {
	filefolder string
	port       string
	upMap      *api.ChanMap
	downMap    *api.ChanMap
}

func (c *Configuration) loadSettings() {
	//TODO: load server settings from somewhere, ex. port number
	c.filefolder = "deadropfiles"
	c.port = ":9090"
	c.upMap = api.InitChanMap()
	c.downMap = api.InitChanMap()
}

var validPath = regexp.MustCompile("^/(create|upload|download|finalize)")

func makeHandler(f func(http.ResponseWriter, *http.Request, *Configuration), conf *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*") //TODO: List of allowed server via config file
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			fmt.Println("invalid path")
			return
		}
		fmt.Println("method:", r.Method)
		f(w, r, conf)
	}
}

func InitServer() *Configuration {
	//TODO: check/start database

	conf := new(Configuration)
	conf.loadSettings()

	//Check if folder "deadropfiles" exist
	if _, err := os.Stat(conf.filefolder); os.IsNotExist(err) {
		err = os.Mkdir(conf.filefolder, 0700) //Borde det vara 0700?
		fmt.Printf("Creating folder %s\n", conf.filefolder)
		if err != nil {
			log.Fatal("Could not create file directory %s\n", err)
		}
	} else {
		fmt.Printf("Folder exists %s\n", conf.filefolder)
	}

	return conf
}

func StartServer(conf *Configuration) {
	// TODO: fix /create, /finalize and /
	//http.HandleFunc("/", makeHandler(upload, conf))
	http.HandleFunc("/create", makeHandler(create, conf))
	http.HandleFunc("/upload", makeHandler(upload, conf))
	http.HandleFunc("/finalize", makeHandler(finalize, conf))
	http.HandleFunc("/download/", makeHandler(download, conf))

	err := http.ListenAndServe(conf.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
