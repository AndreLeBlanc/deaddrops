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

func generateToken() string {
	token := md5.New()
	t := time.Now()
	io.WriteString(token, t.String())
	return hex.EncodeToString(token.Sum(nil))
}

func create(w http.ResponseWriter, r *http.Request, conf *Configuration) {
	if r.Method != "GET" {
		fmt.Println("Create: Invalid request")
		return
	}

	stringToken := generateToken()
	c := make(chan string)
	api.AppendChan(conf.upMap, stringToken, c)

	reply, err := json.Marshal(stringToken)
	if err != nil {
		fmt.Println("Failed token json encoding")
		return
	}

	go api.DummySupervisor2(stringToken, c, conf.upMap)
	fmt.Fprintf(w, string(reply))
}
