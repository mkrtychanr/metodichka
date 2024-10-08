package logger

import (
	"log"
	"os"
	"sync"
)

var (
	initializer    sync.Once
	loggerInstance *log.Logger
)

func GetLogger() *log.Logger {
	initializer.Do(func() {
		loggerInstance = log.New(os.Stdout, "service: ", log.Ldate|log.Ltime)
	})

	return loggerInstance
}
