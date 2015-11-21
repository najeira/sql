package sql

import (
	"bytes"
	"testing"

	log "github.com/najeira/goutils/logv"
	"strings"
)

func TestLggerLevel(t *testing.T) {
	var buf bytes.Buffer
	lgr := log.NewLogger()
	lgr.SetOutput(&buf)
	SetLogger(lgr)

	if lgr != logger {
		t.Errorf("SetLogger: expected %d, got %d", lgr, logger)
	}

	lgr.SetLevel(logDebug)
	if !logv(logDebug) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logInfo) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logWarn) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logErr) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logFatal) {
		t.Errorf("logv: expected true, got false")
	}

	lgr.SetLevel(logInfo)
	if logv(logDebug) {
		t.Errorf("logv: expected false, got true")
	}
	if !logv(logInfo) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logWarn) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logErr) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logFatal) {
		t.Errorf("logv: expected true, got false")
	}

	lgr.SetLevel(logWarn)
	if logv(logDebug) {
		t.Errorf("logv: expected false, got true")
	}
	if logv(logInfo) {
		t.Errorf("logv: expected false, got true")
	}
	if !logv(logWarn) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logErr) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logFatal) {
		t.Errorf("logv: expected true, got false")
	}

	lgr.SetLevel(logErr)
	if logv(logDebug) {
		t.Errorf("logv: expected false, got true")
	}
	if logv(logInfo) {
		t.Errorf("logv: expected false, got true")
	}
	if logv(logWarn) {
		t.Errorf("logv: expected false, got true")
	}
	if !logv(logErr) {
		t.Errorf("logv: expected true, got false")
	}
	if !logv(logFatal) {
		t.Errorf("logv: expected true, got false")
	}

	lgr.SetLevel(logFatal)
	if logv(logDebug) {
		t.Errorf("logv: expected false, got true")
	}
	if logv(logInfo) {
		t.Errorf("logv: expected false, got true")
	}
	if logv(logWarn) {
		t.Errorf("logv: expected false, got true")
	}
	if logv(logErr) {
		t.Errorf("logv: expected false, got true")
	}
	if !logv(logFatal) {
		t.Errorf("logv: expected true, got false")
	}
}

func TestLggerPrint(t *testing.T) {
	var buf bytes.Buffer
	lgr := log.NewLogger()
	lgr.SetOutput(&buf)
	SetLogger(lgr)

	logln("hoge")

	s := buf.String()
	if !strings.Contains(s, "hoge") {
		t.Errorf("logln: expected hoge, got %s", s)
	}
}
