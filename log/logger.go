package log

import (
	"os"
	"github.com/apsdehal/go-logger"
)

var Mylog *logger.Logger

func NewLogger() {
	var (
		err     error
		logFile *os.File
	)
	logFile, err = os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("open log file error.")
	}
	Mylog, err = logger.New("log", 1, logFile)
	if err != nil {
		panic("log init error.")
	}
}
