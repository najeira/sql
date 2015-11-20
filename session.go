package sql

import (
	"database/sql"
	"sync"
)

var (
	sessionPool sync.Pool
)

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

func (s *Session) Close() error {
	// do not close valuesPool at tx session.
	// it will be closed root session.
	if s.tx == nil {
		s.valuesPool.Close()
	}
	s.valuesPool = nil

	s.db = nil
	s.tx = nil

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

func (s *Session) Query(scn Scanner, q string, args ...interface{}) ([]Row, error) {
	return sqlQuery(s.querier(), scn, q, args...)
}

func (s *Session) QueryAsync(scn Scanner, q string, args ...interface{}) chan QueryResult {
	return sqlQueryAsync(s.querier(), scn, q, args...)
}

func (s *Session) Exec(q string, args ...interface{}) (int64, int64, error) {
	return sqlExec(s.executor(), q, args...)
}

func (s *Session) ExecAsync(q string, args ...interface{}) chan ExecResult {
	return sqlExecAsync(s.executor(), q, args...)
}

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

func (s *Session) Commit() error {
	if s.tx == nil {
		return nil // not in tx
	}
	return s.tx.Commit()
}

func (s *Session) Rollback() error {
	if s.tx == nil {
		return nil // not in tx
	}
	return s.tx.Rollback()
}

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
