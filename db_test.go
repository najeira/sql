package sql_test

import (
	"context"
	"testing"

	"github.com/mattn/go-gimei"
	"github.com/najeira/sql"
)

const (
	createTable = "CREATE TABLE IF NOT EXISTS `user` (" +
		"  `id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"  `name` text," +
		"  PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
)

type user struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func open() (*sql.DB, error) {
	return sql.Open(sql.Config{
		User:            "sqltest",
		Passwd:          "testsql",
		ServerName:      "localhost:3306",
		DBName:          "sqltest",
		MaxOpenConns:    0,
		MaxIdleConns:    0,
		ConnMaxLifetime: 0,
	})
}

func TestOpen(t *testing.T) {
	ctx := context.Background()
	db, err := open()
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("nil")
	}

	res, err := db.Exec(ctx, createTable)
	if err != nil {
		t.Error(err)
	} else if res == nil {
		t.Error("nil")
	}
}

func TestQueryer(t *testing.T) {
	ctx := context.Background()
	db, err := open()
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("nil")
	}
	defer db.Close()

	name := gimei.NewName()

	var id int64
	t.Run("insert", func(t *testing.T) {
		q := "insert into `user` (name) values (?)"
		res, err := db.Exec(ctx, q, name.String())
		if err != nil {
			t.Fatal(err)
		}

		id_, err := res.LastInsertId()
		if err != nil {
			t.Error(err)
		}
		id = id_
	})

	t.Run("get", func(t *testing.T) {
		q := "select id, name from `user` where id = ?"
		var row user
		if err := db.Get(ctx, &row, q, id); err != nil {
			t.Fatal(err)
		}

		if row.ID != id {
			t.Error("invalid id", row.ID)
		}
		if row.Name != name.String() {
			t.Error("invalid name", row.Name)
		}
	})

	t.Run("update", func(t *testing.T) {
		q := "update `user` set name = ? where id = ?"
		res, err := db.Exec(ctx, q, name.Hiragana(), id)
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Error(err)
		} else if n != 1 {
			t.Error("invalid RowsAffected")
		}
	})

	t.Run("select", func(t *testing.T) {
		q := "select id, name from `user` where id = ?"
		var rows []*user
		if err := db.Select(ctx, &rows, q, id); err != nil {
			t.Fatal(err)
		}

		if len(rows) != 1 {
			t.Error("invalid rows")
		} else {
			row := rows[0]
			if row.ID != id {
				t.Error("invalid id", row.ID)
			}
			if row.Name != name.Hiragana() {
				t.Error("invalid name", row.Name)
			}
		}
	})

	t.Run("delete", func(t *testing.T) {
		q := "delete from `user` where id = ?"
		res, err := db.Exec(ctx, q, id)
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Error(err)
		} else if n != 1 {
			t.Error("invalid RowsAffected")
		}
	})
}

func TestHooksSelect(t *testing.T) {
	ctx := context.Background()
	db, err := open()
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("nil")
	}
	defer db.Close()

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

	q := "select id, name from `user`"
	var rows []*user
	if err := db.Select(ctx, &rows, q); err != nil {
		t.Fatal(err)
	}

	if q != pre {
		t.Error(pre, q)
	}
	if q != post {
		t.Error(pre, q)
	}
}

func TestMapper(t *testing.T) {
	ctx := context.Background()
	db, err := open()
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("nil")
	}
	defer db.Close()

	name := gimei.NewName()

	q := "insert into `user` (name) values (?)"
	res, err := db.Exec(ctx, q, name.String())
	if err != nil {
		t.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.Error(err)
	}

	db.Mapper("mytag")

	q = "select id, name from `user` where id = ?"
	var row struct {
		ID   int64  `mytag:"id"`
		Name string `mytag:"name"`
	}
	if err := db.Get(ctx, &row, q, id); err != nil {
		t.Fatal(err)
	}

	if row.ID != id {
		t.Error("invalid id", row.ID)
	}
	if row.Name != name.String() {
		t.Error("invalid name", row.Name)
	}
}
