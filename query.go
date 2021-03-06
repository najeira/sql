package sql

import (
	"database/sql"
	"time"
)

type Result struct {
	LastInsertId int64
	RowsAffected int64
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

	tracef("%s %v", q, args)

	metrics.MarkQuery()
	defer timers.Measure(q, time.Now())

	rows, err := svc.Query(q, args...)
	if err != nil {
		errorf("%s", err)
		return nil, err
	}

	return getRowsForSqlRows(rows)
}

func sqlExec(svc executor, q string, args ...interface{}) (Result, error) {
	eres := Result{}
	if svc == nil {
		return eres, errSesionClosed
	}

	args = flattenArgs(args)

	tracef("%s %v", q, args)

	metrics.MarkExecute()
	defer timers.Measure(q, time.Now())

	res, err := svc.Exec(q, args...)
	if err != nil {
		errorf("%s", err)
		return eres, err
	}

	i, err := res.LastInsertId()
	if err != nil {
		errorf("%s", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		errorf("%s", err)
	}

	metrics.MarkAffects(n)

	eres.LastInsertId = i
	eres.RowsAffected = n
	return eres, nil
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
