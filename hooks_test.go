package sql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/najeira/sql"
)

func TestHooks(t *testing.T) {
	var preSelect string
	var postSelect string

	hooks := &sql.Hooks{
		PreSelect: func(
			ctx context.Context,
			dest interface{},
			query string,
			args []interface{},
		) (context.Context, error) {
			preSelect = query
			return ctx, nil
		},
		PostSelect: func(
			ctx context.Context,
			dest interface{},
			query string,
			args []interface{},
			err error,
		) {
			postSelect = query
		},
		PreQuery:     nil,
		PostQuery:    nil,
		PreExec:      nil,
		PostExec:     nil,
		PreBegin:     nil,
		PostBegin:    nil,
		PreCommit:    nil,
		PostCommit:   nil,
		PreRollback:  nil,
		PostRollback: nil,
	}

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})
	db.Hooks(hooks)

	q := "SELECT 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
		}).AddRow(1))

	var rows []int64
	err = db.Select(ctx, &rows, q)
	if err != nil {
		t.Fatal(err)
	}

	if preSelect != q {
		t.Error("hooks PreSelect")
	}
	if postSelect != q {
		t.Error("hooks postSelect")
	}
}
