package log

import (
	"os"
	"github.com/apsdehal/go-logger"
	"fmt"
	"kuaidi100_go/system"
)

var Log *logger.Logger

func NewLogger() {

	var (
		err     error
		logFile *os.File
	)
	logFilePath := fmt.Sprintf("%s%s", system.GetCurrentDirectory(), "/logs.txt")
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("open log file error.")
	}
	Log, err = logger.New("log", 1, logFile)
	if err != nil {
		panic("log init error.")
	}
}
