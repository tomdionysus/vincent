package log

import (
	"fmt"
	"strings"
	"sync"
)

// This is a type to provide ConsoleLogger to the system
type ConsoleLogger struct {
	LogLevel int

	mutex *sync.Mutex
}

// The function creates a New ConsoleLogger with the loglevel supplied
func NewConsoleLogger(logLevel string) *ConsoleLogger {
	logger := &ConsoleLogger{
		LogLevel: parseLogLevel(strings.ToLower(strings.Trim(logLevel, " "))),
		mutex:    &sync.Mutex{},
	}
	if logger.LogLevel == LOG_LEVEL_UNKNOWN {
		logger.Warn("Cannot parse log level '%s', assuming debug", logLevel)
		logger.LogLevel = LOG_LEVEL_DEBUG
	}
	return logger
}

func (clg *ConsoleLogger) GetLogLevel() int         { return clg.LogLevel }
func (clg *ConsoleLogger) SetLogLevel(loglevel int) { clg.LogLevel = loglevel }

// Logs a Raw message (-----) with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Raw(message string, args ...interface{}) {
	clg.printLog("-----", message, args...)
}

// Logs a FATAL message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Fatal(message string, args ...interface{}) {
	clg.printLog("FATAL", message, args...)
}

// Logs an ERROR message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Error(message string, args ...interface{}) {
	clg.printLog("ERROR", message, args...)
}

// Logs a WARN message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Warn(message string, args ...interface{}) {
	if clg.LogLevel > LOG_LEVEL_WARN {
		return
	}
	clg.printLog("WARN ", message, args...)
}

// Logs an INFO message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Info(message string, args ...interface{}) {
	if clg.LogLevel > LOG_LEVEL_INFO {
		return
	}
	clg.printLog("INFO ", message, args...)
}

// Logs a DEBUG message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Debug(message string, args ...interface{}) {
	if clg.LogLevel > LOG_LEVEL_DEBUG {
		return
	}
	clg.printLog("DEBUG", message, args...)
}

func (clg *ConsoleLogger) printLog(level string, message string, args ...interface{}) {
	clg.mutex.Lock()
	defer clg.mutex.Unlock()

	fmt.Printf("%s [%s] ", GetTimeUTCString(), level)
	fmt.Printf(message, args...)
	fmt.Print("\n")
}
