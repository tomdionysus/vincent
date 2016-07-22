package log

import (
	"fmt"
)

// MockLogger is only used in testing, and simply stores the last log output.
type MockLogger struct {
	LastLog  string
	TimeNow  string
	LogLevel int
}

// The function creates a New MockLogger with the loglevel supplied
func NewMockLogger(timeprefix string) *MockLogger {
	return &MockLogger{TimeNow: timeprefix}
}

func (me *MockLogger) GetLogLevel() int         { return me.LogLevel }
func (me *MockLogger) SetLogLevel(loglevel int) { me.LogLevel = loglevel }

// Logs a Raw message (-----) with the specified component, message and Printf-style arguments.
func (me *MockLogger) Raw(message string, args ...interface{}) {
	me.printLog("-----", message, args...)
}

// Logs a FATAL message with the specified component, message and Printf-style arguments.
func (me *MockLogger) Fatal(message string, args ...interface{}) {
	me.printLog("FATAL", message, args...)
}

// Logs an ERROR message with the specified component, message and Printf-style arguments.
func (me *MockLogger) Error(message string, args ...interface{}) {
	me.printLog("ERROR", message, args...)
}

// Logs a WARN message with the specified component, message and Printf-style arguments.
func (me *MockLogger) Warn(message string, args ...interface{}) {
	me.printLog("WARN ", message, args...)
}

// Logs an INFO message with the specified component, message and Printf-style arguments.
func (me *MockLogger) Info(message string, args ...interface{}) {
	me.printLog("INFO ", message, args...)
}

// Logs a DEBUG message with the specified component, message and Printf-style arguments.
func (me *MockLogger) Debug(message string, args ...interface{}) {
	me.printLog("DEBUG", message, args...)
}

func (me *MockLogger) printLog(level string, message string, args ...interface{}) {
	me.LastLog = fmt.Sprintf("%s [%s] ", me.TimeNow, level)
	me.LastLog += fmt.Sprintf(message, args...)
}
