package server

import (
	"log"
	"os"
)

// colors
const (
	PURPLE = "\033[0;35m"
	RED    = "\033[0;31m"
	BLUE   = "\033[0;34m"
	YELLOW = "\033[0;33m"
	RESET  = "\033[0m"
)

var logger = log.New(os.Stdout, "["+BLUE+"chipku"+RESET+"] ", log.LstdFlags|log.Lmicroseconds)

func LogInfo(format string, a ...interface{}) {
	logger.Printf(PURPLE+"[info] "+RESET+format, a...)
}

func LogDebug(format string, a ...interface{}) {
	logger.Printf(YELLOW+"[debug] "+RESET+format, a...)
}

func LogError(format string, a ...interface{}) {
	logger.Printf(RED+"[error] "+RESET+format, a...)
}
