package sql

import (
	"cmp"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

const (
	defaultMaxOpenConns    = 30
	defaultMaxIdleConns    = 30
	defaultConnMaxLifetime = time.Minute
)

type DB struct {
	db *sqlx.DB
}

func Open(cfg Config) (*DB, error) {
	db, err := sql.Open(cfg.driverName(), cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return New(db, cfg), nil
}

func New(db *sql.DB, cfg Config) *DB {
	dbx := sqlx.NewDb(db, cfg.driverName())
	dbx.SetMaxOpenConns(velueOrDefault(cfg.MaxOpenConns, defaultMaxOpenConns))
	dbx.SetMaxIdleConns(velueOrDefault(cfg.MaxIdleConns, defaultMaxIdleConns))
	dbx.SetConnMaxLifetime(velueOrDefault(cfg.ConnMaxLifetime, defaultConnMaxLifetime))
	return &DB{db: dbx}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Mapper(name string) {
	db.db.Mapper = reflectx.NewMapper(name)
}

func (db *DB) Get(ctx context.Context, dest any, query string, args ...any) error {
	return doGet(db.db, ctx, dest, query, args)
}

func (db *DB) Select(ctx context.Context, dest any, query string, args ...any) error {
	return doSelect(db.db, ctx, dest, query, args)
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return doExec(db.db, ctx, query, args)
}

func (db *DB) NamedExec(ctx context.Context, query string, arg any) (sql.Result, error) {
	return doNamedExec(db.db, ctx, query, arg)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return doQuery(db.db, ctx, query, args)
}

func (db *DB) RunInTx(ctx context.Context, p Processor) error {
	txx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	tx := &Tx{tx: txx}
	return tx.runInTx(ctx, p)
}

func (db *DB) Conn(ctx context.Context) (*Conn, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: conn}, nil
}

func velueOrDefault[T cmp.Ordered](value, def T) T {
	var zero T
	if value > zero {
		return value
	}
	return def
}
