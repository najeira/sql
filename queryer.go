package sql

import (
	"context"
	"database/sql"
	"errors"
)

type Processor func(ctx context.Context, db Queryer) error

type Queryer interface {
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExec(ctx context.Context, query string, arg any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type queryer interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type namedExecutor interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

func doGet(q queryer, ctx context.Context, dest any, query string, args []any) error {
	var err error
	err = q.GetContext(ctx, dest, query, args...)
	return err
}

func doSelect(q queryer, ctx context.Context, dest any, query string, args []any) error {
	var err error
	err = q.SelectContext(ctx, dest, query, args...)
	return err
}

func doExec(q queryer, ctx context.Context, query string, args []any) (sql.Result, error) {
	var err error
	result, err := q.ExecContext(ctx, query, args...)
	return result, err
}

func doNamedExec(n namedExecutor, ctx context.Context, query string, arg any) (sql.Result, error) {
	var err error
	result, err := n.NamedExecContext(ctx, query, arg)
	return result, err
}

func doQuery(q queryer, ctx context.Context, query string, args []any) (*sql.Rows, error) {
	var err error
	rows, err := q.QueryContext(ctx, query, args...)
	return rows, err
}

func doCommit(tx *sql.Tx, ctx context.Context) error {
	err := tx.Commit()
	if errors.Is(err, sql.ErrTxDone) {
		err = nil
	}
	return err
}

func doRollback(tx *sql.Tx, ctx context.Context) error {
	err := tx.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		err = nil
	}
	return err
}
