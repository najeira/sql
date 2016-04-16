package sql

import (
	"database/sql"
)

// Session is a database handle.
type Session interface {
	Query(q string, args ...interface{}) (*Rows, error)
	QueryAsync(q string, args ...interface{}) <-chan AsyncRows
	Exec(q string, args ...interface{}) (Result, error)
	ExecAsync(q string, args ...interface{}) <-chan AsyncResult
	Begin() (Session, error)
	Commit() error
	Rollback() error
	RunInTx(f func(Session) error) error
	IsTx() bool
}

// DB is a database handle representing a pool of zero or more
// underlying connections.
type DB struct {
	*sql.DB
}

var _ Session = (*DB)(nil)

// Open opens a database specified by its database driver name and a
// driver-specific data source name.
func Open(driver, dsn string) (*DB, error) {
	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{sqlDB}, nil
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (d *DB) Query(q string, args ...interface{}) (*Rows, error) {
	return sqlQuery(d.DB, q, args...)
}

// QueryAsync executes asynchronously a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (d *DB) QueryAsync(q string, args ...interface{}) <-chan AsyncRows {
	return sqlQueryAsync(d.DB, q, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (d *DB) Exec(q string, args ...interface{}) (Result, error) {
	return sqlExec(d.DB, q, args...)
}

// ExecAsync executes asynchronously a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (d *DB) ExecAsync(q string, args ...interface{}) <-chan AsyncResult {
	return sqlExecAsync(d.DB, q, args...)
}

// Begin starts a transaction.
func (d *DB) Begin() (Session, error) {
	sqlTx, err := d.DB.Begin()
	if err != nil {
		return nil, err
	}
	metrics.MarkBegin()
	return &Tx{tx: sqlTx, cnt: 0}, nil
}

// Commit commits the transaction if the session is for transaction.
func (d *DB) Commit() error {
	return sql.ErrTxDone // not in tx
}

// Rollback aborts the transaction if the session is for transaction.
func (d *DB) Rollback() error {
	return sql.ErrTxDone // not in tx
}

// RunInTx runs the function in a transaction.
func (d *DB) RunInTx(f func(Session) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	tracef("BEGIN")

	ferr := f(tx)
	if ferr != nil {
		if rerr := tx.Rollback(); rerr != nil {
			errorf("ROLLBACK %s", rerr)
		} else {
			tracef("ROLLBACK")
		}
		return ferr
	}

	cerr := tx.Commit()
	if cerr != nil {
		errorf("COMMIT %s", cerr)
		return cerr
	}
	tracef("COMMIT")
	return nil
}

// IsTx returns true if the session for transaction, otherwise false.
func (d *DB) IsTx() bool {
	return false
}
