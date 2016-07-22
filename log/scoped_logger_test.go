package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewScopedLogger(t *testing.T) {
	underinst := NewMockLogger("0001-02-03")

	inst := NewScopedLogger("SCOPE", underinst)

	inst.Raw("TEST MESSAGE")
	assert.Equal(t, "0001-02-03 [-----] SCOPE: TEST MESSAGE", underinst.LastLog)
	inst.Debug("TEST MESSAGE")
	assert.Equal(t, "0001-02-03 [DEBUG] SCOPE: TEST MESSAGE", underinst.LastLog)
	inst.Info("TEST MESSAGE")
	assert.Equal(t, "0001-02-03 [INFO ] SCOPE: TEST MESSAGE", underinst.LastLog)
	inst.Warn("TEST MESSAGE")
	assert.Equal(t, "0001-02-03 [WARN ] SCOPE: TEST MESSAGE", underinst.LastLog)
	inst.Error("TEST MESSAGE")
	assert.Equal(t, "0001-02-03 [ERROR] SCOPE: TEST MESSAGE", underinst.LastLog)
	inst.Fatal("TEST MESSAGE")
	assert.Equal(t, "0001-02-03 [FATAL] SCOPE: TEST MESSAGE", underinst.LastLog)

	inst.SetLogLevel(5)
	assert.Equal(t, 5, underinst.LogLevel)

	underinst.LogLevel = 3
	assert.Equal(t, 3, inst.GetLogLevel())
}
