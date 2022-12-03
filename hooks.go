package sql

import (
	"context"
	"database/sql"
)

type Hooks struct {
	// GetかSelectが呼び出されたときに、その実行の前に呼ばれる
	// 戻り値の `context.Context` は sql.Get sql.Select
	// および `Hooks.PostQuery` に渡される
	// この関数がerrorを返したときはExecは実行されない
	PreSelect func(
		ctx context.Context,
		dest interface{},
		query string,
		args []interface{},
	) (context.Context, error)

	// GetかSelectが実行されたあとに呼び出される
	PostSelect func(
		ctx context.Context,
		dest interface{},
		query string,
		args []interface{},
		err error,
	)

	// Queryが呼び出されたときに、その実行の前に呼ばれる
	// 戻り値の `context.Context` は sql.Query
	// および `Hooks.PostQuery` に渡される
	// この関数がerrorを返したときはExecは実行されない
	PreQuery func(
		ctx context.Context,
		query string,
		args []interface{},
	) (context.Context, error)

	// Queryが実行されたあとに呼び出される
	PostQuery func(
		ctx context.Context,
		query string,
		args []interface{},
		rows *sql.Rows,
		err error,
	)

	// Execが呼び出されたときに、その実行の前に呼ばれる
	// 戻り値の `context.Context` は sql.Exec
	// および `Hooks.PostExec` に渡される
	// この関数がerrorを返したときはExecは実行されない
	PreExec func(
		ctx context.Context,
		query string,
		args []interface{},
	) (context.Context, error)

	// Execが実行されたあとに呼び出される
	PostExec func(
		ctx context.Context,
		query string,
		args []interface{},
		result sql.Result,
		err error,
	)

	PreBegin  func(ctx context.Context) (context.Context, error)
	PostBegin func(ctx context.Context, err error)

	PreCommit  func(ctx context.Context) error
	PostCommit func(ctx context.Context, err error)

	PreRollback  func(ctx context.Context) error
	PostRollback func(ctx context.Context, err error)
}

func (h *Hooks) preSelect(
	ctx context.Context,
	dest interface{},
	query string,
	args []interface{},
) (context.Context, error) {
	if h == nil || h.PreSelect == nil {
		return ctx, nil
	}
	return h.PreSelect(ctx, dest, query, args)
}

func (h *Hooks) postSelect(
	ctx context.Context,
	dest interface{},
	query string,
	args []interface{},
	err error,
) {
	if h == nil || h.PostSelect == nil {
		return
	}
	h.PostSelect(ctx, dest, query, args, err)
}

func (h *Hooks) preQuery(
	ctx context.Context,
	query string,
	args []interface{},
) (context.Context, error) {
	if h == nil || h.PreQuery == nil {
		return ctx, nil
	}
	return h.PreQuery(ctx, query, args)
}

func (h *Hooks) postQuery(
	ctx context.Context,
	query string,
	args []interface{},
	rows *sql.Rows,
	err error,
) {
	if h == nil || h.PostQuery == nil {
		return
	}
	h.PostQuery(ctx, query, args, rows, err)
}

func (h *Hooks) preExec(
	ctx context.Context,
	query string,
	args []interface{},
) (context.Context, error) {
	if h == nil || h.PreExec == nil {
		return ctx, nil
	}
	return h.PreExec(ctx, query, args)
}

func (h *Hooks) postExec(
	ctx context.Context,
	query string,
	args []interface{},
	result sql.Result,
	err error,
) {
	if h == nil || h.PostExec == nil {
		return
	}
	h.PostExec(ctx, query, args, result, err)
}

func (h *Hooks) preBegin(ctx context.Context) (context.Context, error) {
	if h == nil || h.PreBegin == nil {
		return ctx, nil
	}
	return h.PreBegin(ctx)
}

func (h *Hooks) postBegin(ctx context.Context, err error) {
	if h == nil || h.PostBegin == nil {
		return
	}
	h.PostBegin(ctx, err)
}

func (h *Hooks) preCommit(ctx context.Context) error {
	if h == nil || h.PreCommit == nil {
		return nil
	}
	return h.PreCommit(ctx)
}

func (h *Hooks) postCommit(ctx context.Context, err error) {
	if h == nil || h.PostCommit == nil {
		return
	}
	h.PostCommit(ctx, err)
}

func (h *Hooks) preRollback(ctx context.Context) error {
	if h == nil || h.PreRollback == nil {
		return nil
	}
	return h.PreRollback(ctx)
}

func (h *Hooks) postRollback(ctx context.Context, err error) {
	if h == nil || h.PostRollback == nil {
		return
	}
	h.PostRollback(ctx, err)
}
