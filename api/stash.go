package api


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

//Returns the index of a file in the stash
func (s Stash) FindFileInStash(f StashFile) int {
	for i, a := range s.Files {
		if a.Fname == f.Fname {
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

func (s *Stash) RemoveFile(index int){
	s.Files = append(s.Files[:index],s.Files[(index+1):]...)
}
