package sql

import (
	"database/sql"
	"errors"
)

var (
	errRowsClosed = errors.New("sql: Rows are closed")
)

// Rows is the result of a query.
type Rows struct {
	rows    *sql.Rows
	columns []string
}

func getRowsForSqlRows(sr *sql.Rows) (*Rows, error) {
	columns, err := sr.Columns()
	if err != nil {
		errorf("%s", err)
		return nil, err
	}
	return &Rows{rows: sr, columns: columns}, nil
}

// Close closes the Rows, preventing further enumeration. If Next returns
// false, the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (r *Rows) Close() error {
	r.columns = nil
	if r.rows == nil {
		return nil
	}
	err := r.rows.Close()
	r.rows = nil
	return err
}

// Next prepares the next result row for reading with the Fetch method.  It
// returns true on success, or false if there is no next result row or an error
// happened while preparing it.  Err should be consulted to distinguish between
// the two cases.
//
// Every call to Fetch, even the first one, must be preceded by a call to Next.
func (r *Rows) Next() bool {
	return r.rows.Next()
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (r *Rows) Err() error {
	return r.rows.Err()
}

// Scan copies the columns in the current row into the values pointed
// at by dest.
//
// If an argument has type *[]byte, Fetch saves in that argument a copy
// of the corresponding data. The copy is owned by the caller and can
// be modified and held indefinitely. The copy can be avoided by using
// an argument of type *sql.RawBytes instead; see the documentation for
// sql.RawBytes for restrictions on its use.
//
// If an argument has type *interface{}, Fetch copies the value
// provided by the underlying driver without conversion. If the value
// is of type []byte, a copy is made and the caller owns the result.
func (r *Rows) Fetch(dest ...interface{}) (Row, error) {
	if r.rows == nil {
		return nil, errRowsClosed
	}

	err := r.rows.Scan(dest...)
	if err != nil {
		errorf("%s", err)
		return nil, err
	}

	metrics.MarkRow()

	ret := make(Row)
	for i, column := range r.columns {
		ret[column] = dest[i]
	}
	return ret, nil
}
