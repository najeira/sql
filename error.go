package sql

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrConnDone = sql.ErrConnDone
	ErrNoRows   = sql.ErrNoRows
	ErrTxDone   = sql.ErrTxDone
)

// Error number: 1062
// Symbol: ER_DUP_ENTRY
// SQLSTATE: 23000
// Message: Duplicate entry '%s' for key %d
func ErrDuplicateEntry(err error) bool {
	return IsErr(err, 1062)
}

func MakeErrDuplicateEntry(msg string) *mysql.MySQLError {
	return &mysql.MySQLError{
		Number:  1062,
		Message: msg,
	}
}

// Error number: 1452
// Symbol: ER_NO_REFERENCED_ROW_2
// SQLSTATE: 23000
// Message: Cannot add or update a child row: a foreign key constraint fails (%s)
func ErrForeignKeyConstraint(err error) bool {
	return IsErr(err, 1452)
}

// Error number: 3572
// Symbol: ER_LOCK_NOWAIT
// SQLSTATE: HY000
// Message: Statement aborted because lock(s) could not be acquired immediately and NOWAIT is set.
func ErrLockNoWait(err error) bool {
	return IsErr(err, 3572)
}

func IsErr(err error, code uint16) bool {
	if err == nil {
		return false
	}
	var merr *mysql.MySQLError
	ok := errors.As(err, &merr)
	return ok && merr.Number == code
}
