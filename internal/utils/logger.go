package utils

import (
	"log"
)

func GetLogger() *log.Logger {
	logger := log.Default()
	logger.SetPrefix("[codepulse] ")
	return logger
}
