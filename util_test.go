package sql_test

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"

	"github.com/najeira/sql"
)

func TestErrDuplicateEntry(t *testing.T) {
	var err error
	err = sql.MakeErrDuplicateEntry("message")
	if ret := sql.ErrDuplicateEntry(err); !ret {
		t.Error("got false")
	}

	err = errors.New("message")
	if ret := sql.ErrDuplicateEntry(err); ret {
		t.Error("got true")
	}
}

func TestErrForeignKeyConstraint(t *testing.T) {
	var err error
	err = &mysql.MySQLError{
		Number:  1452,
		Message: "message",
	}
	if ret := sql.ErrForeignKeyConstraint(err); !ret {
		t.Error("got false")
	}

	err = errors.New("message")
	if ret := sql.ErrForeignKeyConstraint(err); ret {
		t.Error("got true")
	}
}

func TestErrLockNoWait(t *testing.T) {
	var err error
	err = &mysql.MySQLError{
		Number:  3572,
		Message: "message",
	}
	if ret := sql.ErrLockNoWait(err); !ret {
		t.Error("got false")
	}

	err = errors.New("message")
	if ret := sql.ErrLockNoWait(err); ret {
		t.Error("got true")
	}
}

func TestStringPtr(t *testing.T) {
	if v := sql.StringPtr("foo"); v == nil {
		t.Error("nil")
	}
	if v := sql.StringPtr(""); v != nil {
		t.Error("not nil")
	}
}

func TestNullStringOf(t *testing.T) {
	if v := sql.NullStringOf("foo"); !v.Valid {
		t.Error(v)
	} else if v.String != "foo" {
		t.Error(v)
	}

	if v := sql.NullStringOf(""); v.Valid {
		t.Error(v)
	} else if v.String != "" {
		t.Error(v)
	}
}

func TestStringValue(t *testing.T) {
	s := "foo"
	if v := sql.StringValue(&s); v != s {
		t.Error(v)
	}
	if v := sql.StringValue(nil); v != "" {
		t.Error(v)
	}
}
