package sql

import (
	"database/sql"
)

// DB is a database handle representing a pool of zero or more
// underlying connections.
type DB struct {
	*sql.DB
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name.
func Open(driver, dsn string) (*DB, error) {
	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{sqlDB}, nil
}

// Session creates a new session to execute queries with.
func (d *DB) Session() *Session {
	return getSession(d.DB, nil, nil)
}
