package server

import (
	"deadrop/api"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func upload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request", 400)
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		http.Error(w, "Received bad file", 400)
		return
	}
	defer file.Close()

	token := r.FormValue("token")

	err = validateToken(token, conf)
	if err != nil {
		http.Error(w, "Invalid token", 400)
		return
	}

	err = createFolder(token, conf)
	if err != nil {
		log.Println("Could not create token folder")
		http.Error(w, "Internal server error", 500)
		return
	}
	filename := parseFilename(handler.Filename)

	reply, err := UpSuperUpload(token, filename, conf)
	if err != nil || reply.HttpCode != http.StatusOK {
		log.Println(reply.Message)
		http.Error(w, reply.Message, reply.HttpCode)
		return
	}

	if reply.HttpCode == http.StatusOK {
		f, err := os.OpenFile(filepath.Join(conf.filefolder, token, filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Fprintf(w, reply.Message)

	} else {
		http.Error(w, reply.Message, reply.HttpCode)
	}
}

func parseFilename(path string) string {
	substr := api.ParseURL(path)
	if len(substr) == 0 {
		log.Println("Failed parsing filename")
		return path
	}
	return substr[len(substr)-1]
}

// TODO: Function body should (maybe?) be integrated into api.ValidateToken
func validateToken(token string, conf *Configuration) error {
	if !api.ValidateToken(token) {
		return errors.New("Invalid token, incorrect format")
	}

	_, ok := api.FindChan(conf.upMap, token)
	if !ok {
		return errors.New("Invalid token, could not find in upMap")
	}

	return nil
}

func createFolder(token string, conf *Configuration) error {
	if _, err := os.Stat(filepath.Join(conf.filefolder, token)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(conf.filefolder, token), 0700)
		if err != nil {
			return err
		}
	}
	return nil
}
