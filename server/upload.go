package server

import (
	"deadrop/api"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func upload(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "POST" {
		fmt.Println("Upload: Invalid request")
		http.Error(w, "Invalid request", 400)
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Received bad file", 400)
		return
	}
	defer file.Close()

	token := r.FormValue("token")
	fmt.Printf("[Upload] token: %s\n", token)

	err = validateToken(token, conf)
	if err != nil {
		fmt.Println("Invalid token")
		http.Error(w, "Invalid token", 400)
	}

	api.ValidateFile( /*file*/ ) // ?

	err = createFolder(token, conf)
	if err != nil {
		fmt.Println("Could not create token folder")
		http.Error(w, "Internal server error", 500)
	}
	filename := parseFilename(handler.Filename)

	reply, err := UpSuperUpload(token, filename, conf)
	if err != nil {
		http.Error(w, reply.Message, reply.HttpCode)
	}
	fmt.Println("Sent filename to channel")
	
	if reply.HttpCode == 200 {
		f, err := os.OpenFile(filepath.Join(conf.filefolder, token, filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Fprintf(w, "%v", handler.Header)

	} else {
		http.Error(w, reply.Message, reply.HttpCode)
	}
}

func parseFilename(path string) string {
	substr := api.ParseURL(path)
	if len(substr) == 0 {
		fmt.Println("Failed parsing filename")
		return path
	}
	return substr[len(substr)-1]
}

// TODO: Function body should (maybe?) be integrated into api.ValidateToken
func validateToken(token string, conf *Configuration) error {
	if !api.ValidateToken(token) {
		return errors.New("Invalid token, incorrect format")
	}
	fmt.Println("Checked that token is valid")

	_, ok := api.FindChan(conf.upMap, token)
	if !ok {
		return errors.New("Invalid token, could not find in upMap")
	}
	fmt.Println("Checked that channel exist")

	return nil
}

func createFolder(token string, conf *Configuration) error {
	if _, err := os.Stat(filepath.Join(conf.filefolder, token)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(conf.filefolder, token), 0700)
		fmt.Println("Creating new token folder")
		if err != nil {
			return err
		}
	}
	return nil
}
