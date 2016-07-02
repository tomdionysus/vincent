package log

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestNewConsoleLogger(t *testing.T) {
  inst := NewConsoleLogger("error")

  assert.NotNil(t, inst)
  assert.Equal(t, LOG_LEVEL_ERROR, inst.LogLevel)
} 

func TestNewConsoleLoggerDefaultLevel(t *testing.T) {
  inst := NewConsoleLogger("sljdvbaobscia")

  assert.NotNil(t, inst)
  assert.Equal(t, LOG_LEVEL_DEBUG, inst.LogLevel)
} 

func TestLoggerGetLogLevel(t *testing.T) {
  inst := NewConsoleLogger("debug")

  assert.Equal(t, LOG_LEVEL_DEBUG, inst.GetLogLevel())  
}

func TestLoggerSetLogLevel(t *testing.T) {
  inst := NewConsoleLogger("debug")

  inst.SetLogLevel(LOG_LEVEL_WARN)
  assert.Equal(t, LOG_LEVEL_WARN, inst.LogLevel)  
}


func TestLoggerRaw(t *testing.T) {
  inst := NewConsoleLogger("debug")
  inst.Raw("component","message")
}

func TestLoggerFatal(t *testing.T) {
  inst := NewConsoleLogger("debug")
  inst.Fatal("component","message")
}

func TestLoggerError(t *testing.T) {
  inst := NewConsoleLogger("debug")
  inst.Error("component","message")
}

func TestLoggerWarn(t *testing.T) {
  inst := NewConsoleLogger("debug")
  inst.Warn("component","message")
}

func TestLoggerWarnLevel(t *testing.T) {
  inst := NewConsoleLogger("error")
  inst.Warn("component","message")
}

func TestLoggerInfo(t *testing.T) {
  inst := NewConsoleLogger("debug")
  inst.Info("component","message")
}

func TestLoggerInfoLevel(t *testing.T) {
  inst := NewConsoleLogger("error")
  inst.Info("component","message")
}

func TestLoggerDebug(t *testing.T) {
  inst := NewConsoleLogger("debug")
  inst.Debug("component","message")
}

func TestLoggerDebugLevel(t *testing.T) {
  inst := NewConsoleLogger("error")
  inst.Debug("component","message")
}