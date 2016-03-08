package sql

import (
	"database/sql"
	"errors"
	"sync"
)

var (
	rowsPool        sync.Pool
	disableRowsPool bool

	errRowsClosed   = errors.New("sql: Rows are closed")
	errScannerIsNil = errors.New("sql: Scanner is nil")
)

// Rows is the result of a query.
type Rows struct {
	rows    *sql.Rows
	columns []string
}

func getRowsForSqlRows(r *sql.Rows) (*Rows, error) {
	columns, err := r.Columns()
	if err != nil {
		errorf("%s", err)
		return nil, err
	}
	return getRowsForSqlRowsAndColumns(r, columns), nil
}

func getRowsForSqlRowsAndColumns(rows *sql.Rows, columns []string) *Rows {
	var r *Rows
	if !disableRowsPool {
		poolCounter.Inc(1)
		if v := rowsPool.Get(); v != nil {
			r = v.(*Rows)
		}
	}
	if r == nil {
		newMeter.Mark(1)
		r = &Rows{}
	}
	r.rows = rows
	r.columns = columns
	return r
}

// Close closes the Rows, preventing further enumeration. If Fetch returns
// nil, the Rows are closed automatically. Close is idempotent.
func (r *Rows) Close() error {
	if r.columns == nil {
		return nil
	}
	r.columns = nil

	var err error = nil
	if r.rows != nil {
		err = r.rows.Close()
		r.rows = nil
	}

	if !disableRowsPool {
		rowsPool.Put(r)
		poolCounter.Dec(1)
	}
	return err
}

// FetchOne fetchs the next row.
// It returns a Row on success, or nil if there is no next result row.
// It returns the error, if any, that was encountered during iteration.
func (r *Rows) FetchOne(scn Scanner) (Row, error) {
	if scn == nil {
		return nil, errScannerIsNil
	}
	if r.rows == nil {
		return nil, errRowsClosed
	}

	if !r.rows.Next() {
		return nil, r.rows.Err()
	}

	values, err := scn(r.rows.Scan)
	if err != nil {
		errorf("%s", err)
		return nil, err
	}

	metrics.MarkRow()

	ret := make(Row)
	for i, column := range r.columns {
		ret[column] = values[i]
	}
	return ret, nil
}

// FetchAll fetchs all the rows and close the Rows.
// It returns the error, if any, that was encountered during iteration.
func (r *Rows) Fetch(scn Scanner) ([]Row, error) {
	if scn == nil {
		return nil, errScannerIsNil
	}

	defer r.Close()
	rets := make([]Row, 0)
	for {
		row, err := r.FetchOne(scn)
		if err != nil {
			return nil, err
		}
		if row == nil { // done
			return rets, nil
		}
		rets = append(rets, row)
	}
	panic("unreachable")
}
