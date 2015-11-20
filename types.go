package sql

type Scan func(...interface{}) error

type Scanner func(Scan) ([]interface{}, error)

type DB interface {
	Query(Scanner, string, ...interface{}) ([]Row, error)
	QueryAsync(Scanner, string, ...interface{}) chan QueryResult
	Exec(string, ...interface{}) (int64, int64, error)
	ExecAsync(string, ...interface{}) chan ExecResult
	RunInTx(func(DB) error) error
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
