package tools

import (
	"fmt"
	"os"
	"time"
)

var LogFile *os.File

func OpenLog() (*os.File, error) {
	t := time.Now()
	filename := fmt.Sprintf("%02d-%02d-%d.log", t.Day(), int(t.Month()), t.Year())
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteLog(format string, args ...interface{}) {
	t := time.Now()
	timeStr := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d: ", t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute(), t.Second())
	_, err := fmt.Fprintf(LogFile, timeStr+format+"\n", args...)
	if err != nil {
		return
	}
}
