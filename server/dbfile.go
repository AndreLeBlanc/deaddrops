package server

import (
	"deadrop/api"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func writeJsonFile(s api.Stash, conf *Configuration) bool {
	filename := filepath.Join(conf.filefolder, s.Token, "stash.json")
	j, _ := json.Marshal(s)
	err := ioutil.WriteFile(filename, []byte(j), 0644)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func readJsonFile(token string, conf *Configuration) (api.Stash, error) {
	filename := filepath.Join(conf.filefolder, token, "stash.json")
	f, err := os.Open(filename)
	defer f.Close()
	var s api.Stash
	if err != nil {
		log.Println(err)
		return s, err
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&s)
	if err != nil {
		log.Println(err)
		return s, err
	}
	return s, err
}

func updateJsonFile(s api.Stash, conf *Configuration) bool {
	filename := filepath.Join(conf.filefolder, s.Token, "stash.json")
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
		return false
	}

	j, _ := json.Marshal(s)
	err = ioutil.WriteFile(filename, []byte(j), 0644)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
