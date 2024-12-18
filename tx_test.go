package sql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDB_RunInTx(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := New(d, Config{})

	mock.ExpectBegin()

	q := "SELECT 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
		}).AddRow(1))

	mock.ExpectCommit()

	err = db.RunInTx(ctx, func(ctx context.Context, db Queryer) error {
		var id int64
		err := db.Get(ctx, &id, q)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_RunInTxWithError(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := New(d, Config{})

	mock.ExpectBegin()

	q := "SELECT 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
		}).AddRow(1))

	mock.ExpectRollback()

	err = db.RunInTx(ctx, func(ctx context.Context, db Queryer) error {
		var id int64
		if err := db.Get(ctx, &id, q); err != nil {
			return err
		}
		return errors.New("test")
	})
	if err == nil {
		t.Error("no error")
	}
}

func TestDB_RunInTxWithPanic(t *testing.T) {
	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	ctx := context.Background()
	db := New(d, Config{})

	mock.ExpectBegin()

	q := "SELECT 1"
	mock.ExpectQuery(q).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
		}).AddRow(1))

	mock.ExpectRollback()

	err = db.RunInTx(ctx, func(ctx context.Context, db Queryer) error {
		var id int64
		err := db.Get(ctx, &id, q)
		if err == nil {
			panic("panic")
		}
		return err
	})
	if err == nil {
		t.Error("no error")
	}
}
