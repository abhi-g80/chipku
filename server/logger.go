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

// InfoLog info logger
var InfoLog = log.New(os.Stdout, "["+BLUE+"INFO "+RESET+"] ", log.LstdFlags|log.Lmicroseconds)

// DebugLog debug logger
var DebugLog = log.New(os.Stdout, "["+YELLOW+"DEBUG"+RESET+"] ", log.LstdFlags|log.Lmicroseconds)

// ErrorLog error logger
var ErrorLog = log.New(os.Stdout, "["+RED+"ERROR"+RESET+"] ", log.LstdFlags|log.Lmicroseconds)
