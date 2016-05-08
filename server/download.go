package server

import (
	"deadrop/api"
	"fmt"
	"net/http"
)

func download(w http.ResponseWriter, r *http.Request, cm *api.ChanMap) {
	fmt.Println("method:", r.Method)
	if r.Method != "GET" {
		// Invalid request
		return
	}

	const FileRoot = "root"
	token := r.FormValue("token")
	filename := r.FormValue("filename")
	filepath := FileRoot + "/" + token + "/" + filename

	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-Disposition", "attachment; filename='"+filename+"'")
	http.ServeFile(w, r, filepath)
}
