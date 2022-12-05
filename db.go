package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

const (
	defaultMaxOpenConns    = 30
	defaultMaxIdleConns    = 30
	defaultConnMaxLifetime = time.Second * 60

	defaultDriverName      = "mysql"
)

var (
	_ Queryer = (*DB)(nil)
)

type DB struct {
	db    *sqlx.DB
	hooks *Hooks
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
	dbx.SetMaxOpenConns(maxOpenConns(cfg.MaxOpenConns))
	dbx.SetMaxIdleConns(maxIdleConns(cfg.MaxIdleConns))
	dbx.SetConnMaxLifetime(connMaxLifetime(cfg.ConnMaxLifetime))
	return &DB{
		db: dbx,
	}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Mapper(name string) {
	db.db.Mapper = reflectx.NewMapper(name)
}

func (db *DB) Hooks(hooks *Hooks) {
	db.hooks = hooks
}

func (db *DB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return doGet(db.db, db.hooks, ctx, dest, query, args)
}

func (db *DB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return doSelect(db.db, db.hooks, ctx, dest, query, args)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return doExec(db.db, db.hooks, ctx, query, args)
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return doQuery(db.db, db.hooks, ctx, query, args)
}

func (db *DB) RunInTx(ctx context.Context, p Processor) error {
	var err error
	if ctx, err = db.hooks.preBegin(ctx); err != nil {
		return err
	}

	txx, err := db.db.BeginTxx(ctx, nil)
	db.hooks.postBegin(ctx, err)
	if err != nil {
		return err
	}

	tx := &Tx{
		tx:    txx,
		hooks: db.hooks,
	}
	return tx.runInTx(ctx, p)
}

func (db *DB) InTx() bool {
	return false
}

func (db *DB) Conn(ctx context.Context) (*Conn, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn:  conn,
		hooks: db.hooks,
	}, nil
}

func maxOpenConns(value int) int {
	if value > 0 {
		return value
	}
	return defaultMaxOpenConns
}

func maxIdleConns(value int) int {
	if value > 0 {
		return value
	}
	return defaultMaxIdleConns
}

func connMaxLifetime(value time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return defaultConnMaxLifetime
}
