package sql

import (
	"database/sql"
	"time"
)

// Scan copies the columns in the current row into the arguments.
type Scan func(...interface{}) error

// Scanner returns the columns in the current row.
type Scanner func(Scan) ([]interface{}, error)

type AsyncRows struct {
	Rows *Rows
	Err  error
}

type Result struct {
	LastInsertId int64
	RowsAffected int64
}

type AsyncResult struct {
	LastInsertId int64
	RowsAffected int64
	Err          error
}

type querier interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}

type executor interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func sqlQuery(svc querier, q string, args ...interface{}) (*Rows, error) {
	if svc == nil {
		return nil, errSesionClosed
	}

	args = flattenArgs(args)

	if logv(logTrace) {
		logf("%s %v", q, args)
	}

	metrics.MarkQueries(1)
	defer timers.Measure(q, time.Now())

	rows, err := svc.Query(q, args...)
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
		return nil, err
	}

	return newRows(rows)
}

func sqlQueryAsync(svc querier, q string, args ...interface{}) <-chan AsyncRows {
	ch := make(chan AsyncRows, 1)
	go func() {
		rows, err := sqlQuery(svc, q, args...)
		ch <- AsyncRows{Rows: rows, Err: err}
	}()
	return ch
}

func sqlExec(svc executor, q string, args ...interface{}) (Result, error) {
	eres := Result{}
	if svc == nil {
		return eres, errSesionClosed
	}

	args = flattenArgs(args)

	if logv(logTrace) {
		logf("%s %v", q, args)
	}

	metrics.MarkExecutes(1)
	defer timers.Measure(q, time.Now())

	res, err := svc.Exec(q, args...)
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
		return eres, err
	}

	i, err := res.LastInsertId()
	if err != nil {
		if logv(logWarn) {
			logln(err)
		}
	}

	n, err := res.RowsAffected()
	if err != nil {
		if logv(logWarn) {
			logln(err)
		}
	}

	metrics.MarkAffects(int(n))

	eres.LastInsertId = i
	eres.RowsAffected = n
	return eres, nil
}

func sqlExecAsync(svc executor, q string, args ...interface{}) <-chan AsyncResult {
	ch := make(chan AsyncResult, 1)
	go func() {
		res, err := sqlExec(svc, q, args...)
		ch <- AsyncResult{
			LastInsertId: res.LastInsertId,
			RowsAffected: res.RowsAffected,
			Err:          err,
		}
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
