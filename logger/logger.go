package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{
		Info:  log.New(io.MultiWriter(os.Stdout, f), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(io.MultiWriter(os.Stderr, f), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}
