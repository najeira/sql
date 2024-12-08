package sql

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestErrDuplicateEntry(t *testing.T) {
	var err error
	err = MakeErrDuplicateEntry("message")
	if ret := ErrDuplicateEntry(err); !ret {
		t.Error("got false")
	}

	err = errors.New("message")
	if ret := ErrDuplicateEntry(err); ret {
		t.Error("got true")
	}
}

func TestErrForeignKeyConstraint(t *testing.T) {
	var err error
	err = &mysql.MySQLError{
		Number:  1452,
		Message: "message",
	}
	if ret := ErrForeignKeyConstraint(err); !ret {
		t.Error("got false")
	}

	err = errors.New("message")
	if ret := ErrForeignKeyConstraint(err); ret {
		t.Error("got true")
	}
}

func TestErrLockNoWait(t *testing.T) {
	var err error
	err = &mysql.MySQLError{
		Number:  3572,
		Message: "message",
	}
	if ret := ErrLockNoWait(err); !ret {
		t.Error("got false")
	}

	err = errors.New("message")
	if ret := ErrLockNoWait(err); ret {
		t.Error("got true")
	}
}
