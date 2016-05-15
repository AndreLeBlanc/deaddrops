package api

import (
	"fmt"
	"time"
	"net/http"
)

/*
type StashFile struct {
	Fname    string
	Size     int
	Type     string
	Download int
}

type Stash struct {
	Token    string
	Lifetime int
	Files    []StashFile
}*/

func DummySupervisor1(token string, c chan SuperChan, cm *ChanMap) {
	select {
	case fname := <-c:
		fmt.Printf("received filename: %s\n", fname)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
	}
}

func DummySupervisor2(token string, c chan SuperChan, cm *ChanMap) {
	fmt.Printf("Upload supervisor %s up and running\n", token)
	s := Stash{token, 0, []StashFile{}}
	fmt.Printf("created local stash: %+v", s)
	loop := true
	for loop {
		select {
		case incoming := <-c:
			fmt.Printf("received filename: %+v\n", incoming)
			if s.Token == incoming.Meta.Token {
				s.Files = append(s.Files,incoming.Meta.Files...)
				rc := incoming.C
				fmt.Printf("received filename: %+v\n", s)
				rc <- HttpReplyChan{s, "", http.StatusOK}
			}

		}
	}
}
