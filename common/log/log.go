package log

import (
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func Info(origin, function, msg interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"line":     numLine,
		"file":     srcFile,
		"origin":   origin,
		"function": function,
	}).Info(msg)
}

func Warning(origin, function, msg interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"line":     numLine,
		"file":     srcFile,
		"origin":   origin,
		"function": function,
	}).Warning(msg)
}

func Error(origin, function string, err interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"line":     numLine,
		"file":     srcFile,
		"origin":   origin,
		"function": function,
	}).Error(err)
}

func Debug(origin, function string, value interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"line":     numLine,
		"file":     srcFile,
		"origin":   origin,
		"function": function,
	}).Debug(value)
}

func Fatal(origin, function string, value interface{}) {
	_, path, numLine, _ := runtime.Caller(1)
	srcFile := filepath.Base(path)
	log.WithFields(log.Fields{
		"line":     numLine,
		"file":     srcFile,
		"origin":   origin,
		"function": function,
	}).Fatal(value)
}
