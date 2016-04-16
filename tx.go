package sql

import (
	"database/sql"
	"errors"
)

var (
	errSesionClosed = errors.New("sql: Session is closed")
)

// Tx is a database handle.
type Tx struct {
	tx  *sql.Tx
	cnt int
}

var _ Session = (*Tx)(nil)

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (t *Tx) Query(q string, args ...interface{}) (*Rows, error) {
	return sqlQuery(t.tx, q, args...)
}

// QueryAsync executes asynchronously a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (t *Tx) QueryAsync(q string, args ...interface{}) <-chan AsyncRows {
	return sqlQueryAsync(t.tx, q, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (t *Tx) Exec(q string, args ...interface{}) (Result, error) {
	return sqlExec(t.tx, q, args...)
}

// ExecAsync executes asynchronously a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (t *Tx) ExecAsync(q string, args ...interface{}) <-chan AsyncResult {
	return sqlExecAsync(t.tx, q, args...)
}

// Begin starts a transaction.
func (t *Tx) Begin() (Session, error) {
	// return self because already in transaction
	t.cnt += 1
	return t, nil
}

// Commit commits the transaction if the session is for transaction.
func (t *Tx) Commit() error {
	if t.cnt > 0 {
		t.cnt -= 1
		return nil
	}

	err := t.tx.Commit()
	if err != nil {
		metrics.MarkCommit()
	}
	return err
}

// Rollback aborts the transaction if the session is for transaction.
func (t *Tx) Rollback() error {
	if t.cnt > 0 {
		t.cnt -= 1
		return nil
	}

	err := t.tx.Rollback()
	if err != nil {
		metrics.MarkRollback()
	}
	return err
}

// RunInTx runs the function in a transaction.
func (t *Tx) RunInTx(f func(Session) error) error {
	return f(t)
}

// IsTx returns true if the session for transaction, otherwise false.
func (t *Tx) IsTx() bool {
	return true
}
