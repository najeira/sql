package sql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/najeira/sql"
)

func TestRowsAffected(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})
	defer db.Close()

	t.Run("nil", func(t *testing.T) {
		n := sql.RowsAffected(nil)
		if n != 0 {
			t.Error("nil", n)
		}
	})

	t.Run("sql.Rows", func(t *testing.T) {
		q := "SELECT 1"
		mock.ExpectQuery(q).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
			}).AddRow(1))
		rows, err := db.Query(ctx, q)
		if err != nil {
			t.Fatal(err)
		}

		if n := sql.RowsAffected(rows); n != 0 {
			t.Error("invalid RowsAffected")
		}
	})

	t.Run("sql.Result", func(t *testing.T) {
		const rows = 45
		q := "SELECT 1"
		mock.ExpectExec(q).
			WillReturnResult(
				sqlmock.NewResult(0, rows),
			)
		res, err := db.Exec(ctx, q)
		if err != nil {
			t.Fatal(err)
		}

		if n := sql.RowsAffected(res); n != rows {
			t.Error("invalid RowsAffected")
		}
	})

	t.Run("slice", func(t *testing.T) {
		dest := []string{"foo", "bar"}
		if n := sql.RowsAffected(dest); n != 2 {
			t.Error("invalid RowsAffected")
		}
		if n := sql.RowsAffected(&dest); n != 2 {
			t.Error("invalid RowsAffected")
		}
	})

	t.Run("array", func(t *testing.T) {
		dest := [2]string{"foo", "bar"}
		if n := sql.RowsAffected(dest); n != 2 {
			t.Error("invalid RowsAffected")
		}
		if n := sql.RowsAffected(&dest); n != 2 {
			t.Error("invalid RowsAffected")
		}
	})

	t.Run("string", func(t *testing.T) {
		var dest string
		if n := sql.RowsAffected(dest); n != 1 {
			t.Error("invalid RowsAffected")
		}
		if n := sql.RowsAffected(&dest); n != 1 {
			t.Error("invalid RowsAffected")
		}
	})

	t.Run("bytes", func(t *testing.T) {
		var dest []byte
		if n := sql.RowsAffected(dest); n != 1 {
			t.Error("invalid RowsAffected")
		}
		if n := sql.RowsAffected(&dest); n != 1 {
			t.Error("invalid RowsAffected")
		}
	})
}
