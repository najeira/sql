package sql

import (
	"database/sql"
)

func Open(driver, dsn string) (*DB, error) {
	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{sqlDB}, nil
}

type DB struct {
	*sql.DB
}

func (d *DB) Session() *Session {
	return getSession(d.DB, nil, nil)
}
