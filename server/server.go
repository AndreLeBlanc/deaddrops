package server

import (
	"deadrop/api"
	// "deadrop/database"
	// "database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

//TODO temporary config fake struct
type Configuration struct {
	filefolder string
	port       string
	upMap      *api.ChanMap
	downMap    *api.ChanMap
	uptimeout  time.Duration
	dntimeout  time.Duration
	reqtimeout time.Duration
	logfile    *os.File
	// dbConn     *sql.DB
}

func (c *Configuration) loadSettings() {
	//TODO: load server settings from somewhere, ex. port number
	c.filefolder = "deadropfiles"
	c.port = ":9090"
	c.upMap = api.InitChanMap()
	c.downMap = api.InitChanMap()
	c.uptimeout = 300 //5 min upload timeout
	c.dntimeout = 600 //10 min download timeout
	c.reqtimeout = 1
	// c.dbConn = database.Init()
	f, err := initLog("logfile")
	if err != nil {
		panic("Failed to open logfile")
	}
	c.logfile = f
	c.readStashes()
}

func (c *Configuration) readStashes() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println(err)
		log.Fatal("Error when reading files in deadropsfiles/")
	}
	for _, f := range files {
		sc := make(chan api.SuperChan)
		api.AppendChan(c.downMap, f.Name(), sc)
	}
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
			log.Printf("Invalid URL path %s\n", r.URL.Path)
			http.Error(w, "Invalid URL", 400)
			return
		}
		// fmt.Println("method:", r.Method)
		f(w, r, conf)
	}
}

func InitServer() *Configuration {
	conf := new(Configuration)
	conf.loadSettings()

	if _, err := os.Stat(conf.filefolder); os.IsNotExist(err) {
		err = os.Mkdir(conf.filefolder, 0700)
		if err != nil {
			log.Fatal("Could not create file directory %s\n", err)
		}
	}

	return conf
}

func StartServer(conf *Configuration) {
	// defer database.Close(conf.dbConn)
	defer conf.logfile.Close()
	http.HandleFunc("/create", makeHandler(create, conf))
	http.HandleFunc("/upload", makeHandler(upload, conf))
	http.HandleFunc("/finalize", makeHandler(finalize, conf))
	http.HandleFunc("/download/", makeHandler(download, conf))

	err := http.ListenAndServe(conf.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
