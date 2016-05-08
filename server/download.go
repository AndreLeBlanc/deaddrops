package server

import (
	"fmt"
	"net/http"
	"path/filepath"
)


func download(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	fmt.Println("method:", r.Method)
	if r.Method != "GET" {
		// Invalid request
		return
	}

	token := r.FormValue("token")
	filename := r.FormValue("filename")
	filepath := filepath.Join(conf.Filefolder(), token, filename)

	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Disposition", "attachment; filename='"+filename+"'")
	http.ServeFile(w, r, filepath)
}
