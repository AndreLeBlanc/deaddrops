package api

import (
	"fmt"
	"net/http"
	"time"
)

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
	s := NewEmptyStash()
	s.Token = token
	fmt.Printf("created local stash: %+v", s)
	loop := true
	for loop {
		select {
		case incoming := <-c:
			fmt.Printf("received filename: %+v\n", incoming)
			rc := incoming.C
			if s.Token == incoming.Meta.Token {
				s.Files = append(s.Files, incoming.Meta.Files...)
				fmt.Printf("received filename: %+v\n", s)
				rc <- HttpReplyChan{s, "", http.StatusOK}
			} else {
				rc <- HttpReplyChan{s, "Internal token error", http.StatusInternalServerError}
			}

		}
	}
}
