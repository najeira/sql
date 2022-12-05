package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	_ Queryer = (*Tx)(nil)
)

type Tx struct {
	tx    *sqlx.Tx
	hooks *Hooks
}

func (tx *Tx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return doGet(tx.tx, tx.hooks, ctx, dest, query, args)
}

func (tx *Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return doSelect(tx.tx, tx.hooks, ctx, dest, query, args)
}

func (tx *Tx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return doExec(tx.tx, tx.hooks, ctx, query, args)
}

func (tx *Tx) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return doQuery(tx.tx, tx.hooks, ctx, query, args)
}

func (tx *Tx) InTx() bool {
	return true
}

// DBのRunInTxでBeginされているので
// TxのRunInTxは渡された関数をそのまま実行するだけでよい
func (tx *Tx) RunInTx(ctx context.Context, p Processor) error {
	return p(ctx, tx)
}

// DBのRunInTxから呼び出される
func (tx *Tx) runInTx(ctx context.Context, p Processor) (err error) {
	defer func() {
		// panicはエラーに変換する
		if p := recover(); p != nil {
			if perr, ok := p.(error); ok {
				err = perr
			} else {
				err = fmt.Errorf("%v", p)
			}
		}

		if err != nil {
			// err時はロールバックする
			// ロールバックの失敗は回復できないのでそのまま進む
			// セッションが切れるとロールバックされる
			// Rollbackの結果ではなくもとのerrorを返す
			_ = doRollback(tx.tx.Tx, tx.hooks, ctx)
		} else {
			err = doCommit(tx.tx.Tx, tx.hooks, ctx)
		}
	}()

	err = p(ctx, tx)
	return
}
