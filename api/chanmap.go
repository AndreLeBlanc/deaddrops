package api

import (
	"sync"
)

type SuperChan struct {
	Meta Stash
	C    chan HttpReplyChan
}

type HttpReplyChan struct {
	Meta     Stash
	Message  string
	HttpCode int
}

type ChanMap struct {
	m   map[string]chan SuperChan
	sync.RWMutex
}

// Initialize a new empty ChanMap.
func InitChanMap() *ChanMap {
	return &ChanMap{m: make(map[string]chan SuperChan)}
}

// Append a channel to the ChanMap with the token string as key.
func AppendChan(cm *ChanMap, token string, c chan SuperChan) {
	if _, ok := FindChan(cm, token); ok {
		return
	}
	cm.Lock()
	cm.m[token] = c
	defer cm.Unlock()

}

// Get the channel with the corresponding token string as key.
func FindChan(cm *ChanMap, token string) (chan SuperChan, bool) {
	cm.RLock()
	defer cm.RUnlock()
	c, ok := cm.m[token]
	return c, ok
}

// Deletes a channel from the chanmap.
// Returns true if succesfull
func DeleteChan(cm *ChanMap, token string) bool {
	if _, ok := FindChan(cm, token); !ok {
		return false
	}
	cm.Lock()
	delete(cm.m, token)
	cm.Unlock()
	return true
}

func LenChan(cm *ChanMap) int {
	return len(cm.m)
}
