package sql

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
)

var (
	_ Queryer = (*Tx)(nil)
)

type Tx struct {
	tx *sqlx.Tx
}

func (tx *Tx) Get(ctx context.Context, dest any, query string, args ...any) error {
	return doGet(tx.tx, ctx, dest, query, args)
}

func (tx *Tx) Select(ctx context.Context, dest any, query string, args ...any) error {
	return doSelect(tx.tx, ctx, dest, query, args)
}

func (tx *Tx) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return doExec(tx.tx, ctx, query, args)
}

func (tx *Tx) NamedExec(ctx context.Context, query string, arg any) (sql.Result, error) {
	return doNamedExec(tx.tx, ctx, query, arg)
}

func (tx *Tx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return doQuery(tx.tx, ctx, query, args)
}

func (tx *Tx) runInTx(ctx context.Context, p Processor) (err error) {
	defer func() {
		// panicはエラーに変換する
		if p := recover(); p != nil {
			if perr, ok := p.(error); ok {
				err = fmt.Errorf("%w\n%s", perr, debug.Stack())
			} else {
				err = fmt.Errorf("%v\n%s", p, debug.Stack())
			}
		}

		if err != nil {
			// err時はロールバックする
			// ロールバックの失敗は回復できないのでそのまま進む
			// セッションが切れるとロールバックされる
			// Rollbackの結果ではなくもとのerrorを返す
			_ = doRollback(tx.tx.Tx, ctx)
		} else {
			err = doCommit(tx.tx.Tx, ctx)
		}
	}()

	err = p(ctx, tx)
	return
}
