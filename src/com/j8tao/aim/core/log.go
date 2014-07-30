package core

import (
	"os"
	"log"
	"fmt"
	"time"
	"strings"
	. "com/j8tao/aim/cfg"
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

func LogInfo(v ...interface{}) {
	log.Printf("Info: %v \n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func LogError(v ...interface {}) {
	log.Printf("Error: %v \n", strings.TrimRight(fmt.Sprintln(v...), "\n"))
}

func getLogFile(t time.Time) string {
	return fmt.Sprintf("%s/goserver%d-%d-%d.log", logPath, t.Year(), t.Month(), t.Day())
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
