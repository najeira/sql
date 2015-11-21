package sql

import (
	"database/sql"
)

// Rows is the result of a query.
type Rows struct {
	rows    *sql.Rows
	columns []string
}

func newRows(r *sql.Rows) (*Rows, error) {
	columns, err := r.Columns()
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
		return nil, err
	}
	return &Rows{rows: r, columns: columns}, nil
}

// Close closes the Rows, preventing further enumeration. If Fetch returns
// nil, the Rows are closed automatically. Close is idempotent.
func (r *Rows) Close() error {
	return r.rows.Close()
}

// FetchOne fetchs the next row.
// It returns a Row on success, or nil if there is no next result row.
// It returns the error, if any, that was encountered during iteration.
func (r *Rows) FetchOne(scn Scanner) (Row, error) {
	if !r.rows.Next() {
		return nil, r.rows.Err()
	}

	values, err := scn(r.rows.Scan)
	if err != nil {
		if logv(logErr) {
			logln(err)
		}
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
