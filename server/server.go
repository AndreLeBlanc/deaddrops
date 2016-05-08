package server

import (
	"deadrop/api"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(upload|download)")

func makeHandler(f func(http.ResponseWriter, *http.Request, *api.ChanMap), cm *api.ChanMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			fmt.Println("invalid path")
			return
		}
		f(w, r, cm)
	}
}

func InitServer() {
	//TODO: check/start database

	//TODO: load server settings from somewhere, ex. port number
	cm := api.InitChanMap()
	http.HandleFunc("/upload", makeHandler(createStash, cm))
	http.HandleFunc("/download", makeHandler(download, cm))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
