package sql

import (
	log "github.com/najeira/goutils/logv"
)

const (
	logTrace = log.Trace
	logDebug = log.Debug
	logInfo  = log.Info
	logWarn  = log.Warn
	logErr   = log.Err
	logFatal = log.Fatal
)

var logger log.Logger

func SetLogger(l log.Logger) {
	logger = l
}

func logv(level int) bool {
	return logger != nil && logger.V(level)
}

func logln(v interface{}) {
	if logger != nil {
		logger.Print(v)
	}
}

func logf(f string, args ...interface{}) {
	if logger != nil {
		logger.Printf(f, args...)
	}
}
