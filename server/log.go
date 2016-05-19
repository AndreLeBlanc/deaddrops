package server

import (
	//"fmt"
	"log"
	"os"
	"io"
)

func InitLog(logpath string) (*os.File, error) {
	os.Remove(logpath)
	f, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	// defer f.Close()

	log.SetOutput(io.Writer(f))
	log.Println("***START OF LOG FILE***")

	return f, nil
}
