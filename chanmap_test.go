package main

import (
	"deadrop/api"
	"sync"
	"testing"
)

func TestBasicCM(t *testing.T) {
	cm := api.InitChanMap()
	c1 := make(chan string)
	c2 := make(chan string)
	api.AppendChan(cm, "hello", c1)

	if _, ok := api.FindChan(cm, "hello"); !ok {
		t.Errorf("Could not find token 'hello' in ChanMap")
	}

	api.AppendChan(cm, "hello", c2)

	if c, _ := api.FindChan(cm, "hello"); c != c1 {
		t.Errorf("Overwrite in ChanMap")
	}

	if _, ok := api.FindChan(cm, "bye"); ok {
		t.Errorf("Found token 'bye' in ChanMap")
	}
}

func TestCMconcurrency(t *testing.T) {
	cm := api.InitChanMap()

	var tokens [10]string
	tokens[0] = "a"
	tokens[1] = "b"
	tokens[2] = "c"
	tokens[3] = "d"
	tokens[4] = "a"
	tokens[5] = "f"
	tokens[6] = "b"
	tokens[7] = "c"
	tokens[8] = "a"
	tokens[9] = "q"

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			api.AppendChan(cm, tokens[i], make(chan string))
		}(i)
	}

	wg.Wait()

	for i := 0; i < 10; i++ {
		if _, ok := api.FindChan(cm, tokens[i]); !ok {
			t.Errorf("Could not find token %s in ChanMap", tokens[i])
		}
	}

	if n := api.LenChan(cm); n != 6 {
		t.Errorf("Map containing more/less elements than it should")
	}
}
