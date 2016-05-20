package server

import (
	"io"
	"log"
	"os"
)

// Removes logpath from disc and creates a new logfile. Caller must close file.
//    f.Close()
func initLog(logpath string) (*os.File, error) {
	os.Remove(logpath)
	f, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(io.Writer(f))
	log.Println("***START OF LOG FILE***")

	return f, nil
}
