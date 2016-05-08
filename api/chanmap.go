package api

import (
	"sync"
)


type ChanMap struct {
	m   map[string]chan string
	mux sync.Mutex
}

// Initialize a new empty ChanMap.
func InitChanMap() *ChanMap {
	return &ChanMap{m: make(map[string]chan string)}
}

// Append a channel to the ChanMap with the token string as key.
func AppendChan(cm *ChanMap, token string, c chan string) {
	cm.mux.Lock()
	defer cm.mux.Unlock()
	if _, ok := FindChan(cm, token); ok {
		return
	}
	cm.m[token] = c
}

// Get the channel with the corresponding token string as key.
func FindChan(cm *ChanMap, token string) (chan string, bool) {
	c, ok := cm.m[token]
	return c, ok
}

func LenChan(cm *ChanMap) int {
	return len(cm.m)
}
