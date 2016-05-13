package database

import (
	"fmt"
	"time"
)

func DummySupervisor1(token string, c chan string, cm *ChanMap) {
	select {
	case fname := <-c:
		fmt.Printf("received filename: %s\n", fname)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
	}
}

func DummySupervisor2(token string, c chan string, cm *ChanMap) {
	fmt.Printf("Upload supervisor %s up and running\n", token)
	loop := true
	for loop {
		select {
		case fname := <-c:
			fmt.Printf("received filename: %s\n", fname)
			//	case <-time.After(time.Second + 100)://TODO decide timeout
			//		fmt.Println("timeout")
		}
	}
}
