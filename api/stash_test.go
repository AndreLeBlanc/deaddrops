package api

import (
	"testing"
)

func TestStashDecrementFile(t *testing.T) {
	files := []StashFile{}
	findFile := StashFile{Fname: "test2", Size: 0, Type: "", Download: 0}
	files = append(files, StashFile{Fname: "test1", Size: 0, Type: "", Download: 1})
	files = append(files, StashFile{Fname: "test2", Size: 0, Type: "", Download: 2})
	files = append(files, StashFile{Fname: "test3", Size: 0, Type: "", Download: 3})
	s := NewEmptyStash()
	s.Files = files
	if d := s.DecrementDownloadCounter(findFile); d == 1 && s.Files[s.FindFileInStash(findFile)].Download != d {
		t.Errorf("Incorrect download counter")
	}
}


func TestRemoveStashFile(t *testing.T) {
	files := []StashFile{}
	files = append(files, StashFile{Fname: "test1", Size: 0, Type: "", Download: 1})
	files = append(files, StashFile{Fname: "test2", Size: 0, Type: "", Download: 2})
	files = append(files, StashFile{Fname: "test3", Size: 0, Type: "", Download: 3})
	s := NewEmptyStash()
	s.Files = files
	s.RemoveFile(0)
	if s.FindFileInStash(StashFile{Fname: "test1", Size: 0, Type: "", Download: 1}) != -1 {
		t.Errorf("First file not removed")
	}
	files = []StashFile{}
	files = append(files, StashFile{Fname: "test1", Size: 0, Type: "", Download: 1})
	files = append(files, StashFile{Fname: "test2", Size: 0, Type: "", Download: 2})
	files = append(files, StashFile{Fname: "test3", Size: 0, Type: "", Download: 3})
	s.Files = files
	s.RemoveFile(1)
	if s.FindFileInStash(StashFile{Fname: "test2", Size: 0, Type: "", Download: 2}) != -1 {
		t.Errorf("Middle file not removed")
	}
	files = []StashFile{}
	files = append(files, StashFile{Fname: "test1", Size: 0, Type: "", Download: 1})
	files = append(files, StashFile{Fname: "test2", Size: 0, Type: "", Download: 2})
	files = append(files, StashFile{Fname: "test3", Size: 0, Type: "", Download: 3})
	s.Files = files
	s.RemoveFile(2)
	if s.FindFileInStash(StashFile{Fname: "test3", Size: 0, Type: "", Download: 3}) != -1 {
		t.Errorf("Last file not removed")
	}
}
