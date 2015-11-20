package sql

import (
	"database/sql"
	"time"

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

type Querier interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}

type Executor interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

type Scan func(...interface{}) error
type Scanner func(Scan) ([]interface{}, error)

type DB interface {
	Query(Scanner, string, ...interface{}) ([]Row, error)
	QueryAsync(Scanner, string, ...interface{}) chan QueryResult
	Exec(string, ...interface{}) (int64, int64, error)
	ExecAsync(string, ...interface{}) chan ExecResult
	RunInTx(func(DB) error) error
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

func sqlQuery(sq Querier, s Scanner, q string, args ...interface{}) ([]Row, error) {
	defer Metrics.Measure(time.Now(), q)

	args = flattenArgs(args)

	if logv(log.Debug) {
		logf("%s %v", q, args)
	}

	Metrics.MarkQueries(1)

	rows, err := sq.Query(q, args...)
	if err != nil {
		if logv(log.Err) {
			logln(err)
		}
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		if logv(log.Err) {
			logln(err)
		}
		return nil, err
	}

	rets := make([]Row, 0)

	for rows.Next() {
		values, err := s(rows.Scan)
		if err != nil {
			if logv(log.Err) {
				logln(err)
			}
			return nil, err
		}

		ret := make(Row)
		for i, column := range columns {
			ret[column] = values[i]
		}
		rets = append(rets, ret)
	}

	Metrics.MarkRows(len(rets))

	if logv(log.Debug) {
		logf("%d rows", len(rets))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rets, nil
}

type QueryResult struct {
	Rows []Row
	Err  error
}

func sqlQueryAsync(sq Querier, s Scanner, q string, args ...interface{}) chan QueryResult {
	ch := make(chan QueryResult)
	go func() {
		rows, err := sqlQuery(sq, s, q, args...)
		ch <- QueryResult{Rows: rows, Err: err}
	}()
	return ch
}

func sqlExec(sq Executor, q string, args ...interface{}) (int64, int64, error) {
	defer Metrics.Measure(time.Now(), q)

	args = flattenArgs(args)

	if logv(log.Debug) {
		logf("%s %v", q, args)
	}

	Metrics.MarkExecutes(1)

	res, err := sq.Exec(q, args...)
	if err != nil {
		if logv(log.Err) {
			logln(err)
		}
		return 0, 0, err
	}

	i, err := res.LastInsertId()
	if err != nil {
		if logv(log.Warn) {
			logln(err)
		}
	}

	n, err := res.RowsAffected()
	if err != nil {
		if logv(log.Warn) {
			logln(err)
		}
	}

	Metrics.MarkAffects(int(n))

	return i, n, nil
}

type ExecResult struct {
	LastInsertId int64
	RowsAffected int64
	Err          error
}

func sqlExecAsync(sq Executor, q string, args ...interface{}) chan ExecResult {
	ch := make(chan ExecResult)
	go func() {
		i, n, err := sqlExec(sq, q, args...)
		ch <- ExecResult{LastInsertId: i, RowsAffected: n, Err: err}
	}()
	return ch
}

func flattenArgs(args ...interface{}) []interface{} {
	rets := make([]interface{}, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case []interface{}:
			rets = append(rets, v...)
		default:
			rets = append(rets, v)
		}
	}
	return rets
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
