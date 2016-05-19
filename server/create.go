package server

import (
	"crypto/md5"
	"deadrop/api"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GenerateToken() string {
	token := md5.New()
	t := time.Now()
	io.WriteString(token, t.String())
	return hex.EncodeToString(token.Sum(nil))
}

func create(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "GET" {
		fmt.Println("Create: Invalid request")
		http.Error(w, "Invalid request", 400)
		return
	}

	stringToken := GenerateToken()
	c := make(chan api.SuperChan)
	api.AppendChan(conf.upMap, stringToken, c)
	jsonToken := struct {
		Token string
	}{
		stringToken,
	}
	reply, err := json.Marshal(jsonToken)
	if err != nil {
		fmt.Println("Failed token json encoding")
		http.Error(w, "Internal server error", 500)
		return
	}
	
	go UpSuper(stringToken, conf)

	//TODO: handle error from JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
