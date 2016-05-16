package api

import (
	"sync"
	"testing"
)

func TestBasicCM(t *testing.T) {
	cm := InitChanMap()
	c1 := make(chan SuperChan)
	c2 := make(chan SuperChan)
	AppendChan(cm, "hello", c1)

	if _, ok := FindChan(cm, "hello"); !ok {
		t.Errorf("Could not find token 'hello' in ChanMap")
	}

	AppendChan(cm, "hello", c2)

	if c, _ := FindChan(cm, "hello"); c != c1 {
		t.Errorf("Overwrite in ChanMap")
	}

	if _, ok := FindChan(cm, "bye"); ok {
		t.Errorf("Found token 'bye' in ChanMap")
	}
}

func TestCMconcurrency(t *testing.T) {
	cm := InitChanMap()

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
			AppendChan(cm, tokens[i], make(chan SuperChan))
		}(i)
	}

	wg.Wait()

	for i := 0; i < 10; i++ {
		if _, ok := FindChan(cm, tokens[i]); !ok {
			t.Errorf("Could not find token %s in ChanMap", tokens[i])
		}
	}

	if n := LenChan(cm); n != 6 {
		t.Errorf("Map containing more/less elements than it should")
	}
}


