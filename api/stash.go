package api

type StashFile struct {
	Id       int
	Fname    string
	Size     int
	Type     string
	Download int
}

type Stash struct {
	Token     string
	StashName string
	Lifetime  int
	Files     []StashFile
}

//Returns the index of a file based on filename in the stash
//Returns -1 if file not found
func (s Stash) FindFileInStash(f StashFile) int {
	for i, a := range s.Files {
		if a.Id == f.Id {
			return i
		}
	}
	return -1
}

//Decrements the download counter of file f in the stash.
//Returns the new counter value
func (s *Stash) DecrementDownloadCounter(f StashFile) int {
	i := s.FindFileInStash(f)
	s.Files[i].Download = s.Files[i].Download - 1
	return s.Files[i].Download
}

//Removes a file from a stash
func (s *Stash) RemoveFile(index int) {
	s.Files = append(s.Files[:index], s.Files[(index+1):]...)
}

//Creates a new empty stash struct with all values initialized
//s := Stash{Token: "", StashName: "", Lifetime: 0, 
//Files: []StashFile{}}
func NewEmptyStash() Stash {
	s := Stash{Token: "", StashName: "", Lifetime: 0, Files: []StashFile{}}
	return s
}

//Creates a new empty stash struct with all values initialized
//StashFile{Id: 0, Fname: "", Size: 0, Type: "", Download: -1}
func NewEmptyStashFile() StashFile {
	f := StashFile{Id: 0, Fname: "", Size: 0, Type: "", Download: -1}
	return f
}
