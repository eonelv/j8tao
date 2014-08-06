package core

import (
	. "com/j8tao/aim/cfg"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type ELogger struct {
}

var (
	logPath = "log"
)

func initLog() {
	var logger ELogger
	log.SetOutput(&logger)
	os.MkdirAll(logPath, os.ModeDir)
}

func output(keyworlds string, v ...interface {}) {
	_, file, line, ok := runtime.Caller(2)
	var shortFile string
	if !ok {
		shortFile = "???"
		line = 0
	} else {
		index := strings.LastIndex(file, "/")
		shortFile = string([]byte(file)[index+1:])
	}
	log.Printf("[%v(%d)] %v: %v \n", shortFile, line, keyworlds, strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func LogInfo(v ...interface{}) {
	output("Info", v...)
}

func LogError(v ...interface{}) {
	output("Error", v...)
}

func getLogFile(t time.Time) string {
	return fmt.Sprintf("%s/goserver%s.log", logPath, t.Format("2006-01-02"))
}

func (this *ELogger) Write(p []byte) (int, error) {
	if IsDebug() {
		os.Stdout.Write(p)
	}
	fileName := getLogFile(time.Now())
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		os.Stdout.Write([]byte(fmt.Sprintf("can not open log file:%s r:%v\n", fileName, err)))
		return 0, err
	}

	return file.Write(p)
}
