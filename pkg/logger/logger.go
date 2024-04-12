package logger

import (
	"log"
	"os"
)

var (
	InfoLog = log.New(os.Stdout, "INFO\t", log.LstdFlags)
	ErrLog  = log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)
)
