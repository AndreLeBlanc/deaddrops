package server

import (
	"deadrop/api"
	//"deadrop/database"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Contains meta data relevant to the system supervisor. The Struct should probably
// be in server.go or a future index file. The fields below are suggestions of what
// to implement. Will probably only implement upload and download supervisor counters.
type system struct {
	maxUpSuper int
	maxDnSuper int
	numUpSuper int
	numDnSuper int
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
		log.Printf("SysSuper invalid option: %d\n", sysOpt.option)
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
		log.Printf("Request timeout: %s\n", token)
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
	replyChannel := make(chan api.HttpReplyChan, 1)
	req := api.SuperChan{finalStash, replyChannel}
	return superRequest(finalStash.Token, req, conf.upMap, conf)
}

// The upload supervisor. Has a stash which it updates after every call and writes
// the stash to the database when finalize is called.
func UpSuper(token string, conf *Configuration) {
	c, ok := api.FindChan(conf.upMap, token)
	if !ok {
		log.Printf("Token disappeared from ChanMap: %s\n", token)
		return
	}
	stash := api.NewEmptyStash()
	stash.Token = token
	for {
		select {
		case superChan := <-c:
			replyChan := superChan.C
			if stash.Token == superChan.Meta.Token {
				if superChan.Meta.Lifetime != 0 {
					// TODO: Validate filenames (optional).
					// err := database.InsertStash(conf.dbConn, &superChan.Meta)
					// database.CheckErr(err) // debug
					success := writeJsonFile(superChan.Meta, conf)
					if !success {
						log.Printf("Failed writing json file: %s\n", token)
					}
					// if err != nil {
					// 	replyChan <- api.HttpReplyChan{superChan.Meta, "Failed to write to database", http.StatusInternalServerError}
					// }
					replyChan <- api.HttpReplyChan{superChan.Meta, "Stash completed", http.StatusOK}
					if api.DeleteChan(conf.upMap, token) {
						SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Stash was already completed", http.StatusNotFound})
					} else {
						log.Fatalf("Could not delete channel from chanmap: %s\n", token)
					}
					return
				}

				// TODO: Validate stash restrictions, like number of files in a stash (optional).
				stash.Files = append(stash.Files, superChan.Meta.Files...)
				replyChan <- api.HttpReplyChan{stash, "", http.StatusOK}
			} else {
				replyChan <- api.HttpReplyChan{stash, "Internal token error", http.StatusInternalServerError}
			}
		case <-time.After(time.Second * conf.uptimeout):
			path := filepath.Join(conf.filefolder, token)
			err := os.RemoveAll(path)
			if err != nil {
				log.Println("Error removing stash from disc: %s\n", token)
			}
			if api.DeleteChan(conf.upMap, token) {
				SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Supervisor timeout", http.StatusRequestTimeout})
			} else {
				log.Fatalf("Could not delete channel from chanmap: %s\n", token)
			}
			log.Printf("UpSuper timeout: %s\n", token)
			return
		}
	}
}

func SuperShutdown(c chan api.SuperChan, reply api.HttpReplyChan) {
	for {
		select {
		case r := <-c:
			replyChan := r.C
			replyChan <- reply

		case <-time.After(time.Second * 2):
			return
		}
	}
}

// Contact a download supervisor to get stash.
func DnSuperStash(token string, conf *Configuration) (*api.HttpReplyChan, error) {
	stash := api.Stash{Token: token, Lifetime: 0, Files: []api.StashFile{}}
	replyChannel := make(chan api.HttpReplyChan, 1)
	req := api.SuperChan{stash, replyChannel}
	_, ok := api.FindChan(conf.downMap, token)
	if !ok {
		api.AppendChan(conf.downMap, token, make(chan api.SuperChan, 1))
		go DnSuper(token, conf)
	}
	return superRequest(token, req, conf.downMap, conf)
}

// Contact a download supervisor to download a file.
func DnSuperDownload(token string, fId string, conf *Configuration) (*api.HttpReplyChan, error) {
	stash := api.NewEmptyStash()
	file := api.NewEmptyStashFile()
	id, _ := strconv.Atoi(fId)
	file.Id = id
	stash.Token = token
	stash.Files = append(stash.Files, file)
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
		log.Printf("Token disappeared from ChanMap: %s\n", token)
		return
	}
	stash, sp := readJsonFile(token, conf)
	if sp != nil {
		log.Printf("Stash does not exist: %s\n", token)
		SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "No such stash", http.StatusNotFound})
		return
	}

	// stash := *sp
	for {
		select {
		case superChan := <-c:
			replyChan := superChan.C
			if stash.Token != superChan.Meta.Token {
				log.Printf("Token mismatch: %s\n", token)
				replyChan <- api.HttpReplyChan{stash, "Internal token error", http.StatusInternalServerError}
			}
			length := len(superChan.Meta.Files)
			if length == 0 {
				replyChan <- api.HttpReplyChan{stash, "Current stash status", http.StatusOK}
			} else if length == 1 {
				reqFile := superChan.Meta.Files[0]
				fileIndex := stash.FindFileInStash(reqFile)

				if fileIndex < 0 {
					replyChan <- api.HttpReplyChan{stash, "No such file in stash", http.StatusNotFound}
				} else if stash.Files[fileIndex].Download == 1 {
					replyChan <- api.HttpReplyChan{stash, "Download file OK", http.StatusResetContent}
					stash.RemoveFile(fileIndex)
					if stash.IsEmpty() {
						go RmStash(token, conf)
						if api.DeleteChan(conf.downMap, token) {
							SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Stash no longer available", http.StatusNotFound})
						} else {
							log.Fatalf("Could not delete channel from chanmap: %s\n", token)
						}
						return
					}
					go updateJsonFile(stash, conf)
				} else {
					stash.DecrementDownloadCounter(reqFile)
					go updateJsonFile(stash, conf)
					replyStash := api.NewEmptyStash()
					replyStash.Files = append(replyStash.Files, stash.Files[fileIndex])
					replyChan <- api.HttpReplyChan{replyStash, "Download file OK", http.StatusOK}
				}
			} else {
				log.Printf("Multiple files requested: %s\n", token)
				replyChan <- api.HttpReplyChan{stash, "Bad file handling", http.StatusMethodNotAllowed}
			}
		case <-time.After(time.Second * conf.dntimeout):
			updateJsonFile(stash, conf)
			if api.DeleteChan(conf.downMap, token) {
				SuperShutdown(c, api.HttpReplyChan{api.NewEmptyStash(), "Supervisor timeout", http.StatusRequestTimeout})
			} else {
				log.Fatalf("Could not delete channel from chanmap: %s\n", token)
			}
			log.Printf("DnSuper timeout: %s\n", token)
			return

		}
	}
}

// Sends a message to the janitor to remove a file from a stash. Both
// in the database and on disc.
func RmFile(token string, fname string, conf *Configuration) {
	err := os.Remove(filepath.Join(conf.filefolder, token, fname))
	if err != nil {
		log.Printf("Failed to remove file: %s/%s\n", token, fname)
	}
}

// Sends a message to the janitor to remove a stash. Both in the
// database and on disc.
func RmStash(token string, conf *Configuration) {
	time.Sleep(time.Second * conf.reqtimeout)
	err := os.RemoveAll(filepath.Join(conf.filefolder, token))
	if err != nil {
		log.Printf("Failed to remove stash: %s/%s\n", token)
	}
}
