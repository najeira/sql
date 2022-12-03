package sql_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/najeira/sql"
)

func TestNew(t *testing.T) {
	d, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	db := sql.New(d, sql.Config{
		User:            "user",
		Passwd:          "password",
		ServerName:      "localhost:3306",
		DBName:          "mydb",
		MaxOpenConns:    0,
		MaxIdleConns:    0,
		ConnMaxLifetime: 0,
	})
	if db == nil {
		t.Error("nil")
	}
}

func TestDB_Get(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})

	q := "SELECT id, name FROM users WHERE id = 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(1, "Taro"))

	var user struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if err := db.Get(ctx, &user, q); err != nil {
		t.Fatal(err)
	}
	if user.ID != 1 {
		t.Error("invalid ID")
	}
	if user.Name != "Taro" {
		t.Error("invalid name")
	}
}

func TestDB_Select(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})

	q := "SELECT id, name FROM users WHERE id = 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(1, "Taro"))

	var users []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	err = db.Select(ctx, &users, q)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal("invalid rows")
	}
	if users[0].ID != 1 {
		t.Error("invalid ID")
	}
	if users[0].Name != "Taro" {
		t.Error("invalid name")
	}
}

func TestDB_Exec(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})

	q := "UPDATE user SET name = ? WHERE id = ?"
	mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(0, 1))

	res, err := db.Exec(ctx, q, "Tarou", 1)
	if err != nil {
		t.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.Error(err)
	}
	if id != 0 {
		t.Error("invalid LastInsertId")
	}

	num, err := res.RowsAffected()
	if err != nil {
		t.Error(err)
	}
	if num != 1 {
		t.Error("invalid RowsAffected")
	}
}

func TestDB_Query(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})

	q := "SELECT id, name FROM users WHERE id = 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(1, "Taro"))

	rows, err := db.Query(ctx, q)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user struct {
			ID   int64  `db:"id"`
			Name string `db:"name"`
		}
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			t.Error(err)
		}
		if user.ID != 1 {
			t.Error("invalid ID")
		}
		if user.Name != "Taro" {
			t.Error("invalid name")
		}
	}
	if err := rows.Err(); err != nil {
		t.Error(err)
	}
}

func TestDB_HooksSelect(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := sql.New(d, sql.Config{})

	var pre string
	var post string
	db.Hooks(&sql.Hooks{
		PreSelect: func(ctx context.Context, dest interface{}, query string, args []interface{}) (context.Context, error) {
			pre = query
			return ctx, nil
		},
		PostSelect: func(ctx context.Context, dest interface{}, query string, args []interface{}, err error) {
			post = query
		},
	})

	q := "SELECT id, name FROM users WHERE id = 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name",
		}).AddRow(1, "Taro"))

	var user struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if err := db.Get(ctx, &user, q); err != nil {
		t.Fatal(err)
	}

	if q != pre {
		t.Error(pre, q)
	}
	if q != post {
		t.Error(pre, q)
	}
}
