package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"deadrop/api"
	"fmt"
	"path/filepath"
)

func writeJsonFile(s api.Stash, conf *Configuration) bool {
	filename := filepath.Join(conf.filefolder, s.Token, "stash.json")
	fmt.Println(filename)
	j, _ := json.Marshal(s)
	err := ioutil.WriteFile(filename,[]byte(j), 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func readJsonFile(token string, conf *Configuration) (api.Stash, error) {
	filename := filepath.Join(conf.filefolder, token, "stash.json")
	fmt.Println(filename)
	f, err := os.Open(filename)
	defer f.Close()
	var s api.Stash
	if err != nil {
		fmt.Printf("the error is ", err)
		return s, err
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&s)
	if err != nil {
		fmt.Printf("the error is ", err)
		return s, err
	}
	return s, err
}
