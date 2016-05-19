package server

import (
	"deadrop/api"
	//"deadrop/database"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"os"
)

// Contains meta data relevant to the system supervisor. The Struct should probably
// be in server.go or a future index file. The fields below are suggestions of what
// to implement. Will probably only implement upload and download supervisor counters.
type system struct {
	maxUpSuper int // Maximum number of upload supervisors
	maxDnSuper int // Maximum number of download supervisors
	numUpSuper int // Number of active upload supervisors
	numDnSuper int // Number of active download supervisors
}

// Codes for the SysOption struct's option field.
const (
	ADDUP int = 101
	ADDDN int = 102
	DECUP int = 103
	DECDN int = 104
)

// A struct to be sent to the system supervisor before making an upload or download
// supervisor. Where respChan is where the response from the system supervisor is
// sent back to.
type SysOption struct {
	option   int
	respChan chan bool
}

// The channel to send requests for confirmation to, should probably be in the
// Configuration struct in server.go.
var SysChan = make(chan SysOption)

// Initialize the system.
func InitSys(maxUp int, maxDn int) *system {
	return &system{maxUp, maxDn, 0, 0}
}

// Send an upload request to system supervisor. Returns true if ok.
func SysUpRequest(sysChan chan SysOption) error {
	return sysRequest(sysChan, ADDDN)
}

// Send a download request to system supervisor. Returns true if ok.
func SysDnRequest(sysChan chan SysOption) error {
	return sysRequest(sysChan, ADDDN)
}

// Notify system supervisor that an upload supervisor has terminated.
func SysUpDone(sysChan chan SysOption) error {
	return sysRequest(sysChan, DECUP)
}

// Notify system supervisor that a download supervisor has terminated.
func SysDnDone(sysChan chan SysOption) error {
	return sysRequest(sysChan, DECDN)
}

func sysRequest(sysChan chan SysOption, option int) error {
	respChan := make(chan bool)
	sysChan <- SysOption{option, respChan}
	select {
	case resp := <-respChan:
		if resp {
			return nil
		}
		return errors.New("Upload request denied")
	}
}

// The system supervisor.
func SysSuper(sys *system, sysChan chan SysOption) {
	for true {
		select {
		case sysOpt := <-sysChan:
			fmt.Println("[SysSuper] received SysOption")
			err := sysHandler(sys, sysOpt)
			switch err {
			case nil:
				sysOpt.respChan <- true
			default:
				sysOpt.respChan <- false
			}
		}
	}
}

// Handles the option sent to the system supervisor.
func sysHandler(sys *system, sysOpt SysOption) error {
	switch sysOpt.option {
	case ADDUP:
		if sys.numUpSuper >= sys.maxUpSuper {
			return errors.New("No more uploads allowed")
		}
		sys.numUpSuper++
	case ADDDN:
		if sys.numDnSuper >= sys.maxDnSuper {
			return errors.New("No more downloads allowed")
		}
		sys.numDnSuper++
	case DECUP:
		sys.numUpSuper--
	case DECDN:
		sys.numDnSuper--
	default:
		fmt.Println("[SysSuper] invalid SysOption")
		return errors.New("Invalid SysOption")
	}

	return nil
}

func superRequest(token string, req api.SuperChan, cm *api.ChanMap, conf *Configuration) (*api.HttpReplyChan, error) {
	respChan := req.C
	c, ok := api.FindChan(cm, token)
	if !ok {
		return &api.HttpReplyChan{Message: "Invalid token", HttpCode: http.StatusServiceUnavailable}, errors.New("Invalid token")
	}
	c <- req
	select {
	case resp := <-respChan:
		return &resp, nil
	case <-time.After(time.Second * conf.reqtimeout):
		return &api.HttpReplyChan{Message: "timeout", HttpCode: http.StatusServiceUnavailable}, nil
	}
}

// Contact an upload supervisor to add a file to a stash.
func UpSuperUpload(token string, fname string, conf *Configuration) (*api.HttpReplyChan, error) {
	newFile := api.NewEmptyStashFile()
	newFile.Fname = fname
	stash := api.NewEmptyStash()
	stash.Token = token
	stash.Files = append(stash.Files, newFile)
	replyChannel := make(chan api.HttpReplyChan, 1)
	req := api.SuperChan{stash, replyChannel}
	return superRequest(token, req, conf.upMap, conf)
}

// Contact an upload supervisor to finalize its stash and end the upload session.
func UpSuperFinalize(finalStash api.Stash, conf *Configuration) (*api.HttpReplyChan, error) {
	replyChannel := make(chan api.HttpReplyChan)
	req := api.SuperChan{finalStash, replyChannel}
	return superRequest(finalStash.Token, req, conf.upMap, conf)
}

// The upload supervisor. Has a stash which it updates after every call and writes
// the stash to the database when finalize is called.
func UpSuper(token string, conf *Configuration) {
	// Create a stash.
	c, ok := api.FindChan(conf.upMap, token)
	if !ok {
		fmt.Println("[UpSuper]: Invalid token")
		return
	}
	stash := api.NewEmptyStash()
	stash.Token = token
	fmt.Printf("UpSuper %s running\n", token)
	// Update is for every request.
	for {
		select {
		case superChan := <-c:
			//fmt.Printf("received : %+v\n", sc)
			replyChan := superChan.C
			if stash.Token == superChan.Meta.Token {
				if superChan.Meta.Lifetime != 0 {
					// TODO: Validate filenames (optional).
					//err := database.InsertStash(conf.dbConn, &superChan.Meta)
					//database.CheckErr(err) // debug
					success := writeJsonFile(superChan.Meta, conf)
					fmt.Println(success)
					/*
						if err != nil {
							replyChan <- api.HttpReplyChan{superChan.Meta, "Failed to write to database", http.StatusInternalServerError}
						}*/
					replyChan <- api.HttpReplyChan{superChan.Meta, "Stash completed", http.StatusOK}
					//TODO chanmap should be cleaned up
					if api.DeleteChan(conf.upMap, token) {
						SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Stash was already completed", http.StatusNotFound})
					} else {
						panic("Could not delete channel from chanmap!")
					}
					return
				}

				// TODO: Validate stash restrictions, like number of files in a stash (optional).
				stash.Files = append(stash.Files, superChan.Meta.Files...)
				fmt.Printf("received filename: %+v\n", stash)
				replyChan <- api.HttpReplyChan{stash, "", http.StatusOK}
			} else {
				replyChan <- api.HttpReplyChan{stash, "Internal token error", http.StatusInternalServerError}
			}
		case  <- time.After(time.Second*conf.uptimeout):
			path := filepath.Join(conf.filefolder, token)
			err := os.RemoveAll(path)
			if err != nil {
				//logg
			}
			if api.DeleteChan(conf.upMap, token) {
				SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Supervisor timeout", http.StatusRequestTimeout})
			} else {
				panic("Could not delete channel from chanmap!!")
			}
			fmt.Println("Supervisor timeout")
			return
		}
	}
}

func SuperShutdown(c chan api.SuperChan, reply api.HttpReplyChan) {
	for {
		select {
		case r := <- c:
			replyChan := r.C
			replyChan <- reply

		case  <- time.After(time.Second * 2):
			return
		}
	}
}

// Contact a download supervisor to get stash.
func DnSuperStash(token string, conf *Configuration) (*api.HttpReplyChan, error) {
	stash := api.Stash{Token: token, Lifetime: 0, Files: []api.StashFile{}}
	replyChannel := make(chan api.HttpReplyChan)
	req := api.SuperChan{stash, replyChannel}
	_, ok := api.FindChan(conf.downMap, token)
	if !ok {
		api.AppendChan(conf.downMap, token, make(chan api.SuperChan, 1))
		go DnSuper(token, conf)
	}
	return superRequest(token, req, conf.downMap, conf)
}

// Contact a download supervisor to download a file.
func DnSuperDownload(token string, fname string, conf *Configuration) (*api.HttpReplyChan, error) {
	stash := api.Stash{Token: token, Lifetime: 0, Files: append([]api.StashFile{}, api.StashFile{Fname: fname, Size: 0, Type: "", Download: 0})}
	replyChannel := make(chan api.HttpReplyChan, 1)
	req := api.SuperChan{stash, replyChannel}
	return superRequest(token, req, conf.downMap, conf)
}

// The download supervisor. Reads the stash from the database and keeps track of
// stash meta data, like Lifetime and #downloads. Calls the janitor to remove files
// from the database and disc when their #downloads reach 0 and the whole stash when
// Lifetime expires.
func DnSuper(token string, conf *Configuration) {
	//stash := api.NewEmptyStash()
	//sp := database.SelectStash(conf.dbConn, token)
	c, ok := api.FindChan(conf.downMap, token)
	if !ok {
		fmt.Println("[DnSuper]: Invalid token")
		return
	}
	stash, sp := readJsonFile(token, conf)
	if sp != nil {
		fmt.Println("[DnSuper]: Stash does not exist")
		SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "No such stash", http.StatusNotFound})
		return
	}
	fmt.Println("I just found this stash ")
	fmt.Println(stash)
	//	stash := *sp
	//	stash.Token = token // this is just temporary

	fmt.Printf("DnSuper %s running\n", token)
	for {
		select {
		case superChan := <-c:
			replyChan := superChan.C
			if stash.Token != superChan.Meta.Token {
				replyChan <- api.HttpReplyChan{stash, "Internal token error", http.StatusInternalServerError}
			}
			length := len(superChan.Meta.Files)
			if length == 0 {
				replyChan <- api.HttpReplyChan{stash, "Current stash status", http.StatusOK}
			} else if length == 1 {
				reqFile := superChan.Meta.Files[0]
				fileIndex := stash.FindFileInStash(reqFile)
				//fmt.Println(fileIndex)
				fmt.Println(stash)
				if fileIndex < 0 {
					replyChan <- api.HttpReplyChan{stash, "No such file in stash", http.StatusNotFound}
				} else if stash.Files[fileIndex].Download == 0 {
					replyChan <- api.HttpReplyChan{stash, "File no longer available", http.StatusNotFound}
				} else {
					n := stash.DecrementDownloadCounter(reqFile)
					if n == 0 {
						//defer RmFile(stash)
					}
					replyChan <- api.HttpReplyChan{stash, "Download file OK", http.StatusOK}
				}
			} else {
				replyChan <- api.HttpReplyChan{stash, "Bad file handling", http.StatusInternalServerError}
			}
			case <-time.After(time.Second * conf.dntimeout):
			updateJsonFile(stash, conf)
			if api.DeleteChan(conf.downMap, token) {
				SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Supervisor timeout", http.StatusRequestTimeout})
			} else {
				panic("Could not delete channel from chanmap!!")
			}
			fmt.Println("Supervisor timeout")
			return
			
		}
	}
}

// Sends a message to the janitor to remove a file from a stash. Both
// in the database and on disc.
func RmFile(fname string) error {
	return errors.New("bla")
}

// Sends a message to the janitor to remove a stash. Both in the
// database and on disc.
func RmStash(token string) error {
	return errors.New("bla")
}
