package utils

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[CAMPAIGN] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf("ERROR: "+format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf("INFO: "+format, v...)
}
