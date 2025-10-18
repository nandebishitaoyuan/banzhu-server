package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

type simpleLogger struct {
	*log.Logger
}

func New() Logger {
	return &simpleLogger{
		Logger: log.New(os.Stdout, "[GoServer] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *simpleLogger) Info(v ...interface{}) {
	l.Println(v...)
}

func (l *simpleLogger) Infof(format string, v ...interface{}) {
	l.Printf(format, v...)
}
