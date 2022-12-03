package sql_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/najeira/sql"
)

func TestRowsAffected(t *testing.T) {
	rows := sqlmock.NewRows([]string{
		"id", "name",
	}).AddRow(1, "Taro")
	n := sql.RowsAffected(rows)
	if n != 1 {
		t.Error("invalid RowsAffected")
	}

	result := sqlmock.NewResult(123, 45)
	n = sql.RowsAffected(result)
	if n != 45 {
		t.Error("invalid RowsAffected")
	}

	dest := []string{"foo", "bar"}
	n = sql.RowsAffected(dest)
	if n != 2 {
		t.Error("invalid RowsAffected")
	}
}
