package sql_test

import (
	sqld "database/sql"
	"testing"

	"github.com/najeira/sql"
)

func TestErrors(t *testing.T) {
	if sql.ErrNoRows != sqld.ErrNoRows {
		t.Error("ErrNoRows")
	}
	if sql.ErrConnDone != sqld.ErrConnDone {
		t.Error("ErrConnDone")
	}
	if sql.ErrTxDone != sqld.ErrTxDone {
		t.Error("ErrTxDone")
	}
}
