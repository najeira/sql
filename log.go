package sql

type Logger interface {
	V(level int) bool
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func tracef(f string, args ...interface{}) {
	if logger != nil {
		logger.Tracef(f, args...)
	}
}

func debugf(f string, args ...interface{}) {
	if logger != nil {
		logger.Debugf(f, args...)
	}
}

func infof(f string, args ...interface{}) {
	if logger != nil {
		logger.Infof(f, args...)
	}
}

func warnf(f string, args ...interface{}) {
	if logger != nil {
		logger.Warnf(f, args...)
	}
}

func errorf(f string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(f, args...)
	}
}
