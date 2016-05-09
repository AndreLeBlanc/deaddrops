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
	chanMap    *api.ChanMap
}

func (c *Configuration) loadSettings() {
	//TODO: load server settings from somewhere, ex. port number
	c.filefolder = "deadropfiles"
	c.port = ":8080"
	c.chanMap = api.InitChanMap()
}

func (c *Configuration) Filefolder() string {
	return c.filefolder
}

func (c *Configuration) ChanMap() *api.ChanMap {
	return c.chanMap
}

var validPath = regexp.MustCompile("^/(upload|download)")

func makeHandler(f func(http.ResponseWriter, *http.Request, *Configuration), conf *Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			fmt.Println("invalid path")
			return
		}
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
		fmt.Printf("Creating folder %s", conf.filefolder)
		if err != nil {
			log.Fatal("Could not create file directory %s\n", err)
		}
	} else {
		fmt.Printf("Folder exists %s\n", conf.filefolder)
	}

	return conf
}

func StartServer(conf *Configuration) {
	http.HandleFunc("/upload", makeHandler(upload, conf))
	http.HandleFunc("/download", makeHandler(download, conf))

	err := http.ListenAndServe(conf.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
