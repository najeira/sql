package sql

import (
	"testing"

	_ "github.com/proullon/ramsql/driver"
)

func TestDB(t *testing.T) {
	db, err := Open("ramsql", "TestLoadUserAddresses")
	if err != nil {
		t.Errorf("Open: error %v", err)
	}
	if db == nil {
		t.Errorf("Open: nil")
	}
	db.Close()
}

func TestSession(t *testing.T) {
	db, err := Open("ramsql", "TestLoadUserAddresses")
	if err != nil {
		t.Errorf("Open: error %v", err)
	}
	defer db.Close()

	var q string
	var res Result
	var rows *Rows

	q = "CREATE TABLE user (id INT PRIMARY KEY AUTOINCREMENT, name TEXT, age INT);"
	res, err = db.Exec(q)
	if err != nil {
		t.Errorf("Session.Exec: %s error %v", q, err)
	}
	if res.LastInsertId != 0 {
		t.Errorf("Session.Exec: LastInsertId %d", res.LastInsertId)
	}
	if res.RowsAffected != 1 {
		t.Errorf("Session.Exec: RowsAffected %d", res.RowsAffected)
	}

	q = "INSERT INTO user (name, age) VALUES ('Akihabara', 32);"
	res, err = db.Exec(q)
	if err != nil {
		t.Errorf("Session.Exec: %s error %v", q, err)
	}
	if res.LastInsertId == 0 {
		t.Errorf("Session.Exec: LastInsertId %d", res.LastInsertId)
	}
	if res.RowsAffected != 1 {
		t.Errorf("Session.Exec: RowsAffected %d", res.RowsAffected)
	}

	users := map[string]int{
		"Iidabashi":   23,
		"Ueno":        23,
		"Okachimachi": 31,
	}
	for name, age := range users {
		q = "INSERT INTO user (name, age) VALUES (?, ?);"
		res, err = db.Exec(q, name, age)
		if err != nil {
			t.Errorf("Session.Exec: %s error %v", q, err)
		}
	}

	q = "SELECT id, name, age FROM user WHERE age = ?;"
	rows, err = db.Query(q, 32)
	if err != nil {
		t.Errorf("Session.Query: %s error %v", q, err)
	}
	if rows == nil {
		t.Errorf("Session.Query: %s nil", q)
	}
	defer rows.Close()

	if rows.Next() {
		var id NullInt64
		var name NullString
		var age NullInt64
		row, err := rows.Fetch(&id, &name, &age)
		if err != nil {
			t.Errorf("Rows.Fetch: error %v", err)
		}
		if row == nil {
			t.Errorf("Rows.Fetch: nil")
		}
		if row.Int("id") != 1 {
			t.Errorf("Row.String: expected 1, got %d", row.Int("id"))
		}
		if row.String("name") != "Akihabara" {
			t.Errorf("Row.String: expected Akihabara, got %s", row.String("name"))
		}
		if row.Int("age") != 32 {
			t.Errorf("Row.String: expected 32, got %d", row.Int("age"))
		}
	}

	if rows.Next() {
		t.Errorf("Rows.Next: not false")
	}
	
	if err := rows.Err(); err != nil {
		t.Errorf("Rows.Err: error %v", err)
	}
}
