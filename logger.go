package godeconz

import (
	"log"
	"os"
)

// Logger supports trace for maximal information, debug for general debug information and error in case of an error.
type Logger interface {
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type BuiltinLogger struct {
	logger *log.Logger
}

func NewBuiltinLogger() *BuiltinLogger {
	return &BuiltinLogger{logger: log.New(os.Stdout, "", 5)}
}

func (l BuiltinLogger) Debugf(format string, args ...interface{}) {
	l.logger.Print("[Debug]")
	l.logger.Printf(format, args...)
}

func (l BuiltinLogger) Errorf(format string, args ...interface{}) {
	l.logger.Print("[Error]")
	l.logger.Printf(format, args...)
}

func (l BuiltinLogger) Warnf(format string, args ...interface{}) {
	l.logger.Print("[Warn]")
	l.logger.Printf(format, args...)
}

func (l BuiltinLogger) Infof(format string, args ...interface{}) {
	l.logger.Print("[Info]")
	l.logger.Printf(format, args)
}
