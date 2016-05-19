package server

import (
	"crypto/md5"
	"deadrop/api"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func generateToken() string {
	token := md5.New()
	t := time.Now()
	io.WriteString(token, t.String())
	return hex.EncodeToString(token.Sum(nil))
}

func create(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request", 400)
		return
	}

	stringToken := generateToken()
	c := make(chan api.SuperChan)
	api.AppendChan(conf.upMap, stringToken, c)
	jsonToken := struct {
		Token string
	}{
		stringToken,
	}
	reply, err := json.Marshal(jsonToken)
	if err != nil {
		log.Printf("Failed token json encoding: %s\n", stringToken)
		http.Error(w, "Internal server error", 500)
		return
	}

	go UpSuper(stringToken, conf)

	w.Header().Set("Content-Type", "application/json")
	w.Write(reply)
}
