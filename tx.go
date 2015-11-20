package sql

import (
	"database/sql"
)

type Tx struct {
	*sql.Tx
}

var _ DB = (*Tx)(nil)

func (t *Tx) Query(s Scanner, q string, args ...interface{}) ([]Row, error) {
	return sqlQuery(t.Tx, s, q, args...)
}

func (t *Tx) QueryAsync(s Scanner, q string, args ...interface{}) chan QueryResult {
	return sqlQueryAsync(t.Tx, s, q, args...)
}

func (t *Tx) Exec(q string, args ...interface{}) (int64, int64, error) {
	return sqlExec(t.Tx, q, args...)
}

func (t *Tx) ExecAsync(q string, args ...interface{}) chan ExecResult {
	return sqlExecAsync(t.Tx, q, args...)
}

func (t *Tx) RunInTx(f func(DB) error) error {
	return f(t)
}
