package sql

type Scan func(...interface{}) error

type Scanner func(Scan) ([]interface{}, error)

type Valuer interface {
	String() *NullString
	Int64() *NullInt64
	Float64() *NullFloat64
	Bool() *NullBool
}

type QueryResult struct {
	Rows []Row
	Err  error
}

type ExecResult struct {
	LastInsertId int64
	RowsAffected int64
	Err          error
}
