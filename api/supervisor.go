package api

import (
	"fmt"
	"time"
)

func DummySupervisor(token string, c chan string, cm *ChanMap) {
	select {
	case fname := <-c:
		fmt.Println("received filename: %s", fname)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
	}
}
