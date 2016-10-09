package log

type ScopedLogger struct {
	Logger Logger
	Scope  string
}

func NewScopedLogger(name string, logger Logger) *ScopedLogger {
	return &ScopedLogger{
		Scope:  name,
		Logger: logger,
	}
}

func (slg *ScopedLogger) GetLogLevel() int         { return slg.Logger.GetLogLevel() }
func (slg *ScopedLogger) SetLogLevel(loglevel int) { slg.Logger.SetLogLevel(loglevel) }

// Raw logs a Raw message (-----) with the specified message and Printf-style arguments.
func (slg *ScopedLogger) Raw(message string, args ...interface{}) {
	slg.Logger.Raw(slg.Scope+": "+message, args...)
}

// Fatal logs a FATAL message with the specified message and Printf-style arguments.
func (slg *ScopedLogger) Fatal(message string, args ...interface{}) {
	slg.Logger.Fatal(slg.Scope+": "+message, args...)
}

// Error logs an ERROR message with the specified message and Printf-style arguments.
func (slg *ScopedLogger) Error(message string, args ...interface{}) {
	slg.Logger.Error(slg.Scope+": "+message, args...)
}

// Warn logs a WARN message with the specified message and Printf-style arguments.
func (slg *ScopedLogger) Warn(message string, args ...interface{}) {
	slg.Logger.Warn(slg.Scope+": "+message, args...)
}

// Info logs an INFO message with the specified message and Printf-style arguments.
func (slg *ScopedLogger) Info(message string, args ...interface{}) {
	slg.Logger.Info(slg.Scope+": "+message, args...)
}

// Debug logs a DEBUG message with the specified message and Printf-style arguments.
func (slg *ScopedLogger) Debug(message string, args ...interface{}) {
	slg.Logger.Debug(slg.Scope+": "+message, args...)
}
