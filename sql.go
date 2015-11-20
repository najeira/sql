package sql

import (
	"database/sql"
)

func Open(driver, dsn string) (DB, error) {
	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &Db{sqlDB}, nil
}

type Db struct {
	*sql.DB
}

var _ DB = (*Db)(nil)

func (d *Db) Query(s Scanner, q string, args ...interface{}) ([]Row, error) {
	return sqlQuery(d.DB, s, q, args...)
}

func (d *Db) QueryAsync(s Scanner, q string, args ...interface{}) chan QueryResult {
	return sqlQueryAsync(d.DB, s, q, args...)
}

func (d *Db) Exec(q string, args ...interface{}) (int64, int64, error) {
	return sqlExec(d.DB, q, args...)
}

func (d *Db) ExecAsync(q string, args ...interface{}) chan ExecResult {
	return sqlExecAsync(d.DB, q, args...)
}

func (d *Db) Begin() (*Tx, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (d *Db) RunInTx(f func(DB) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	if logv(logDebug) {
		logln("BEGIN")
	}
	err = f(tx)
	if err != nil {
		tx.Rollback()
		if logv(logDebug) {
			logln("ROLLBACK")
		}
	} else {
		tx.Commit()
		if logv(logDebug) {
			logln("COMMIT")
		}
	}
	return err
}
