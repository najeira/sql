package sql

import (
	"database/sql"

	log "github.com/najeira/goutils/logv"
	"github.com/najeira/goutils/maputil"
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
	if logv(log.Debug) {
		logln("BEGIN")
	}
	err = f(tx)
	if err != nil {
		tx.Rollback()
		if logv(log.Debug) {
			logln("ROLLBACK")
		}
	} else {
		tx.Commit()
		if logv(log.Debug) {
			logln("COMMIT")
		}
	}
	return err
}

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

func Placeholders(num int) string {
	size := len("?") + len(",")
	buf := make([]byte, num*size-1)
	for i := 0; i < num; i++ {
		buf[i*size] = []byte("?")[0]
		if i != num {
			buf[i*size+1] = []byte(",")[0]
		}
	}
	return string(buf)
}

func IntsToArgs(nums []int64) []interface{} {
	args := make([]interface{}, len(nums))
	for i, n := range nums {
		args[i] = n
	}
	return args
}

func CollectInts(rows []Row, key string) []int64 {
	rets := make([]int64, len(rows))
	for i, row := range rows {
		ret, _ := maputil.Int(row, key)
		rets[i] = ret
	}
	return rets
}
