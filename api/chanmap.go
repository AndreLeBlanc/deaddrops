package api

import (
	"sync"
)

type SuperChan struct{
	Meta Stash
	C chan HttpReplyChan
}

type HttpReplyChan struct{
	Meta Stash
	Message string
	HttpCode int
}

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
}

type ChanMap struct {
	m   map[string]chan SuperChan
	mux sync.Mutex
}

// Initialize a new empty ChanMap.
func InitChanMap() *ChanMap {
	return &ChanMap{m: make(map[string]chan SuperChan)}
}

// Append a channel to the ChanMap with the token string as key.
func AppendChan(cm *ChanMap, token string, c chan SuperChan) {
	cm.mux.Lock()
	defer cm.mux.Unlock()
	if _, ok := FindChan(cm, token); ok {
		return
	}
	cm.m[token] = c
}

// Get the channel with the corresponding token string as key.
func FindChan(cm *ChanMap, token string) (chan SuperChan, bool) {
	c, ok := cm.m[token]
	return c, ok
}

func LenChan(cm *ChanMap) int {
	return len(cm.m)
}
