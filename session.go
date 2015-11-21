package sql

import (
	"database/sql"
	"sync"
)

var (
	sessionPool sync.Pool
)

// Session is a database handle.
type Session struct {
	*valuesPool

	db *sql.DB
	tx *sql.Tx
}

func getSession(db *sql.DB, tx *sql.Tx, vp *valuesPool) *Session {
	poolCounter.Inc(1)
	var s *Session
	if v := sessionPool.Get(); v != nil {
		s = v.(*Session)
	} else {
		s = &Session{}
	}
	s.db = db
	s.tx = tx
	if vp != nil {
		s.valuesPool = vp
	} else {
		s.valuesPool = getValuesPool()
	}
	return s
}

// Close closes the Session.
func (s *Session) Close() error {
	if s.db == nil && s.tx == nil {
		return nil
	}

	// do not close valuesPool at tx session.
	// it will be closed root session.
	if s.tx == nil {
		s.valuesPool.Close()
	}
	s.valuesPool = nil

	s.db = nil
	s.tx = nil

	// put this Session to the pool.
	sessionPool.Put(s)
	poolCounter.Dec(1)
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
		return s, nil
	}

	sqlTx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	tx := getSession(nil, sqlTx, s.valuesPool)
	return tx, nil
}

// Commit commits the transaction if the session is for transaction.
func (s *Session) Commit() error {
	if s.tx == nil {
		return sql.ErrTxDone // not in tx
	}
	return s.tx.Commit()
}

// Rollback aborts the transaction if the session is for transaction.
func (s *Session) Rollback() error {
	if s.tx == nil {
		return sql.ErrTxDone // not in tx
	}
	return s.tx.Rollback()
}

// RunInTx runs the function in a transaction.
func (s *Session) RunInTx(f func(*Session) error) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}
	if logv(logDebug) {
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
			if logv(logDebug) {
				logln("ROLLBACK")
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		if logv(logDebug) {
			logln("COMMIT")
		}
	}
	return err
}

// IsTx returns true if the session for transaction, otherwise false.
func (s *Session) IsTx() bool {
	return s.tx != nil
}
