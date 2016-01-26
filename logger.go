package main

import (
	"log"
	"os"
)

const (
	logLevelFatal = iota
	logLevelError
	logLevelInfo
	logLevelDebug
)

// AppLogger is an interface that wraps different log level logging methods
type AppLogger interface {
	Fatal(f string)
	Error(f string)
	Info(f string)
	Debug(f string)
}

type applogger struct {
	level int
	log   *log.Logger
}

func newAppLogger(level int) *applogger {
	return &applogger{
		level: level,
		log:   log.New(os.Stderr, "", 0),
	}
}

// Fatal logs fatal level logs
func (l *applogger) Fatal(f string) {
	l.output(logLevelFatal, "FATAL", f)
}

// Error logs error level logs
func (l *applogger) Error(f string) {
	l.output(logLevelError, "ERROR", f)
}

// Info logs info level logs
func (l *applogger) Info(f string) {
	l.output(logLevelInfo, "INFO", f)
}

// Debug logs debug level logs
func (l *applogger) Debug(f string) {
	l.output(logLevelDebug, "DEBUG", f)
}

func (l *applogger) output(level int, prefix string, f string) {
	if level <= l.level {
		l.log.Printf("# [%s] %s", prefix, f)
	}

	if level == logLevelFatal {
		os.Exit(1)
	}
}
