package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Conn struct {
	conn *sqlx.Conn
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

func (conn *Conn) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return doGet(conn.conn, ctx, dest, query, args)
}

func (conn *Conn) Ping(ctx context.Context) error {
	return conn.conn.PingContext(ctx)
}
