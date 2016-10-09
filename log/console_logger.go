package log

import (
	"fmt"
	"strings"
	"sync"
)

// ConsoleLogger is a type to provide ConsoleLogger to the system
type ConsoleLogger struct {
	LogLevel int

	mutex *sync.Mutex
}

// NewConsoleLogger creates a New ConsoleLogger with the loglevel supplied
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

// GetLogLevel gets the current log level
func (clg *ConsoleLogger) GetLogLevel() int { return clg.LogLevel }

// SetLogLevel sets the current log level
func (clg *ConsoleLogger) SetLogLevel(loglevel int) { clg.LogLevel = loglevel }

// Raw logs a Raw message (-----) with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Raw(message string, args ...interface{}) {
	clg.printLog("-----", message, args...)
}

// Fatal logs a FATAL message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Fatal(message string, args ...interface{}) {
	clg.printLog("FATAL", message, args...)
}

// Error logs an ERROR message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Error(message string, args ...interface{}) {
	clg.printLog("ERROR", message, args...)
}

// Warn logs a WARN message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Warn(message string, args ...interface{}) {
	if clg.LogLevel > LOG_LEVEL_WARN {
		return
	}
	clg.printLog("WARN ", message, args...)
}

// Info logs an INFO message with the specified message and Printf-style arguments.
func (clg *ConsoleLogger) Info(message string, args ...interface{}) {
	if clg.LogLevel > LOG_LEVEL_INFO {
		return
	}
	clg.printLog("INFO ", message, args...)
}

// Debug logs a DEBUG message with the specified message and Printf-style arguments.
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
