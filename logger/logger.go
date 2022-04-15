package logger

import (
	"log"
	"path/filepath"
	"runtime"
)

var logger *log.Logger

func SetDefault(l *log.Logger) {
	logger = l
}

func Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok == false {
		file, line = "", -1
	} else {
		file = filepath.Base(file)
	}
	args = append(args, "file:")
	args = append(args, file)
	args = append(args, "line:")
	args = append(args, line)

	logger.Println(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Printf(format, args...)
}
