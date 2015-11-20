package sql

import (
	"database/sql"
	"errors"
	"time"
)

var ErrSessionClosed = errors.New("sql: Session has already been closed")

type querier interface {
	Query(string, ...interface{}) (*sql.Rows, error)
}

type executor interface {
	Exec(string, ...interface{}) (sql.Result, error)
}

func sqlQuery(svc querier, s Scanner, q string, args ...interface{}) ([]Row, error) {
	if svc == nil {
		return nil, ErrSessionClosed
	}

	defer Metrics.Measure(time.Now(), q)

	args = flattenArgs(args)

	if logv(logDebug) {
		logf("%s %v", q, args)
	}

	Metrics.MarkQueries(1)

	rows, err := svc.Query(q, args...)
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
		return nil, err
	}

	rets := make([]Row, 0)

	for rows.Next() {
		values, err := s(rows.Scan)
		if err != nil {
			if logv(logErr) {
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

	if logv(logDebug) {
		logf("%d rows", len(rets))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rets, nil
}

func sqlQueryAsync(svc querier, s Scanner, q string, args ...interface{}) chan QueryResult {
	ch := make(chan QueryResult)
	go func() {
		rows, err := sqlQuery(svc, s, q, args...)
		ch <- QueryResult{Rows: rows, Err: err}
	}()
	return ch
}

func sqlExec(svc executor, q string, args ...interface{}) (int64, int64, error) {
	if svc == nil {
		return 0, 0, ErrSessionClosed
	}

	defer Metrics.Measure(time.Now(), q)

	args = flattenArgs(args)

	if logv(logDebug) {
		logf("%s %v", q, args)
	}

	Metrics.MarkExecutes(1)

	res, err := svc.Exec(q, args...)
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
		return 0, 0, err
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

	Metrics.MarkAffects(int(n))

	return i, n, nil
}

func sqlExecAsync(svc executor, q string, args ...interface{}) chan ExecResult {
	ch := make(chan ExecResult)
	go func() {
		i, n, err := sqlExec(svc, q, args...)
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
