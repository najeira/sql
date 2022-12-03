package sql

import (
	"context"
	"database/sql"
)

type Processor func(ctx context.Context, db Queryer) error

type Queryer interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	RunInTx(ctx context.Context, p Processor) error
	InTx() bool
}

type queryer interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func doGet(q queryer, h *Hooks, ctx context.Context, dest interface{}, query string, args []interface{}) error {
	var err error
	if ctx, err = h.preSelect(ctx, dest, query, args); err != nil {
		return err
	}
	err = q.GetContext(ctx, dest, query, args...)
	metrics.MarkGet(err)
	h.postSelect(ctx, dest, query, args, err)
	return err
}

func doSelect(q queryer, h *Hooks, ctx context.Context, dest interface{}, query string, args []interface{}) error {
	var err error
	if ctx, err = h.preSelect(ctx, dest, query, args); err != nil {
		return err
	}
	err = q.SelectContext(ctx, dest, query, args...)
	metrics.MarkSelect(dest, err)
	h.postSelect(ctx, dest, query, args, err)
	return err
}

// nolint:interfacer
func doExec(q queryer, h *Hooks, ctx context.Context, query string, args []interface{}) (sql.Result, error) {
	var err error
	if ctx, err = h.preExec(ctx, query, args); err != nil {
		return nil, err
	}
	result, err := q.ExecContext(ctx, query, args...)
	metrics.MarkResult(result, err)
	h.postExec(ctx, query, args, result, err)
	return result, err
}

func doQuery(q queryer, h *Hooks, ctx context.Context, query string, args []interface{}) (*sql.Rows, error) {
	var err error
	if ctx, err = h.preQuery(ctx, query, args); err != nil {
		return nil, err
	}
	rows, err := q.QueryContext(ctx, query, args...)
	metrics.MarkQuery()
	h.postQuery(ctx, query, args, rows, err)
	return rows, err
}

func doCommit(tx *sql.Tx, h *Hooks, ctx context.Context) error {
	if err := h.preCommit(ctx); err != nil {
		return err
	}

	err := tx.Commit()
	if err == sql.ErrTxDone {
		err = nil
	}

	metrics.MarkCommit()
	h.postCommit(ctx, err)
	return err
}

func doRollback(tx *sql.Tx, h *Hooks, ctx context.Context) error {
	if err := h.preRollback(ctx); err != nil {
		return err
	}

	err := tx.Rollback()
	if err == sql.ErrTxDone {
		err = nil
	}

	metrics.MarkRollback()
	h.postRollback(ctx, err)
	return err
}
