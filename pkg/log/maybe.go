package log

type maybeLogger struct {
	parent Logger
}

func (m maybeLogger) Info(message string, args ...interface{}) {
	if m.parent != nil {
		m.parent.Info(message,args...)
	}
}

func (m maybeLogger) Error(message string, args ...interface{}) {
	if m.parent != nil {
		m.parent.Error(message, args...)
	}
}

func (m maybeLogger) Debug(message string, args ...interface{}) {
	if m.parent != nil {
		m.parent.Debug(message, args...)
	}
}

func MaybeLogger(logger Logger) Logger {
	return &maybeLogger{logger}
}

