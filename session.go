package sql

import (
	"database/sql"
	"errors"
	"sync"
)

var (
	sessionPool        sync.Pool
	disableSessionPool bool

	errSesionClosed = errors.New("sql: Session is closed")
)

// Session is a database handle.
type Session struct {
	db      *sql.DB
	tx      *sql.Tx
	txCount int
	values  *values
}

func getSession(db *sql.DB, tx *sql.Tx, vp *values) *Session {
	var s *Session = nil
	if !disableSessionPool {
		poolCounter.Inc(1)
		if v := sessionPool.Get(); v != nil {
			s = v.(*Session)
		}
	}
	if s == nil {
		newMeter.Mark(1)
		s = &Session{}
	}

	s.db = db
	s.tx = tx
	if vp != nil {
		s.values = vp
	} else {
		s.values = getValues()
	}
	return s
}

// Close closes the Session.
func (s *Session) Close() error {
	if s.db == nil && s.tx == nil {
		return nil
	}

	// do not close values at tx session.
	// it will be cleared by root session.
	if s.tx == nil {
		s.values.Clear()
	}
	s.values = nil

	s.db = nil
	s.tx = nil
	s.txCount = 0

	// put this Session to the pool.
	if !disableSessionPool {
		sessionPool.Put(s)
		poolCounter.Dec(1)
	}
	return nil
}

func (s *Session) querier() querier {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) executor() executor {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (s *Session) Query(q string, args ...interface{}) (*Rows, error) {
	return sqlQuery(s.querier(), q, args...)
}

// QueryAsync executes asynchronously a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (s *Session) QueryAsync(q string, args ...interface{}) <-chan AsyncRows {
	return sqlQueryAsync(s.querier(), q, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (s *Session) Exec(q string, args ...interface{}) (Result, error) {
	return sqlExec(s.executor(), q, args...)
}

// ExecAsync executes asynchronously a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (s *Session) ExecAsync(q string, args ...interface{}) <-chan AsyncResult {
	return sqlExecAsync(s.executor(), q, args...)
}

// Begin starts a transaction.
func (s *Session) Begin() (*Session, error) {
	// return self if already in transaction
	if s.tx != nil {
		s.txCount += 1
		return s, nil
	}

	if s.db == nil {
		return nil, errSesionClosed
	}

	sqlTx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	metrics.MarkBegin()

	tx := getSession(nil, sqlTx, s.values)
	return tx, nil
}

// Commit commits the transaction if the session is for transaction.
func (s *Session) Commit() error {
	if s.tx == nil {
		return sql.ErrTxDone // not in tx
	}

	if s.txCount > 0 {
		s.txCount -= 1
		return nil
	}

	err := s.tx.Commit()
	if err != nil {
		metrics.MarkCommit()
	}
	return err
}

// Rollback aborts the transaction if the session is for transaction.
func (s *Session) Rollback() error {
	if s.tx == nil {
		return sql.ErrTxDone // not in tx
	}

	if s.txCount > 0 {
		s.txCount -= 1
		return nil
	}

	err := s.tx.Rollback()
	if err != nil {
		metrics.MarkRollback()
	}
	return err
}

// RunInTx runs the function in a transaction.
func (s *Session) RunInTx(f func(*Session) error) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}
	if logv(logTrace) {
		logln("BEGIN")
	}

	err = f(tx)

	if err != nil {
		rerr := tx.Rollback()
		if rerr != nil {
			if logv(logErr) {
				logf("ROLLBACK %v", rerr)
			}
		} else {
			if logv(logTrace) {
				logln("ROLLBACK")
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		if logv(logErr) {
			logf("COMMIT %v", err)
		}
	} else {
		if logv(logTrace) {
			logln("COMMIT")
		}
	}
	return err
}

// IsTx returns true if the session for transaction, otherwise false.
func (s *Session) IsTx() bool {
	return s.tx != nil
}

func (s *Session) String() *NullString {
	if s.values == nil {
		return nil
	}
	return s.values.String()
}

func (s *Session) Int64() *NullInt64 {
	if s.values == nil {
		return nil
	}
	return s.values.Int64()
}

func (s *Session) Float64() *NullFloat64 {
	if s.values == nil {
		return nil
	}
	return s.values.Float64()
}

func (s *Session) Bool() *NullBool {
	if s.values == nil {
		return nil
	}
	return s.values.Bool()
}
