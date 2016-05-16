package server

import (
	"fmt"
	//"time"
	"deadrop/api"
	"errors"
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

// Meta data for an upload supervisor.
type UploadStash struct {
	token string
	// db  *database
	//upChan chan SomeStruct // Should have a response channel field
}

func superRequest(token string, req api.SuperChan, cm *api.ChanMap) (*api.HttpReplyChan, error) {
	respChan := req.C
	c, ok := api.FindChan(cm, token)
	if !ok {
		return nil, errors.New("Invalid token")
	}
	c <- req
	select {
	case resp := <-respChan:
		return &resp, nil
	}
}

// Contact an upload supervisor to add a file to a stash.
func UpSuperUpload(token string, fname string, conf *Configuration) (*api.HttpReplyChan, error) {
	stash := api.Stash{Token: token, Lifetime: 0, Files: append([]api.StashFile{}, api.StashFile{Fname: fname, Size: 0, Type: "", Download: 0})}
	replyChannel := make(chan api.HttpReplyChan)
	req := api.SuperChan{stash, replyChannel}
	return superRequest(token, req, conf.upMap)
}

// Contact an upload supervisor to finalize its stash and end the upload session.
func UpSuperFinalize(finalStash api.Stash, conf *Configuration) (*api.HttpReplyChan, error) {
	replyChannel := make(chan api.HttpReplyChan)
	req := api.SuperChan{finalStash, replyChannel}
	return superRequest(finalStash.Token, req, conf.upMap)
}

// The upload supervisor. Has a stash which it updates after every call and writes
// the stash to the database when finalize is called.
func UpSuper( /*stashUp stashUploadStruct*/ ) {
	// Create a stash.
	// Update is for every request.
	// Write stash to database when it gets a finalize flag.
}

// Meta data for an download supervisor. Will need a response channel to send a json.
type DownloadStash struct {
	token string
	// db   *database
	//dnChan chan SomeStruct // Should have a response channel field
	janitor chan string
}

// Contact an upload supervisor to download a file.
func DnSuperDownload(token string, fname string, conf *Configuration) (*api.HttpReplyChan, error) {
	stash := api.Stash{Token: token, Lifetime: 0, Files: append([]api.StashFile{}, api.StashFile{Fname: fname, Size: 0, Type: "", Download: 0})}
	replyChannel := make(chan api.HttpReplyChan)
	req := api.SuperChan{stash, replyChannel}
	return superRequest(token, req, conf.downMap)
}

// The download supervisor. Reads the stash from the database and keeps track of
// stash meta data, like Lifetime and #downloads. Calls the janitor to remove files
// from the database and disc when their #downloads reach 0 and the whole stash when
// Lifetime expires.
func DnSuper() {
	// Read stash from database.
	// Update it for every request.
}

// Metadata for the cleaning supervisor.
type CleanStash struct {
	option       int
	token        string
	fname        string
	responseChan bool
}

// Constants used for the option field in the CleanStash struct.
const (
	CLEANFILE  = 200
	CLEANSTASH = 201
)

// The channel to send cleaning requests to, should probably be in the
// Configuration struct in server.go.
var CleanChan = make(chan CleanStash)

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

// The janitor. Download supervisor tells the janitor what to do.
func Cleaner() {

}
