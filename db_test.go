package sql

import (
	"testing"

	_ "github.com/proullon/ramsql/driver"
)

func TestDB(t *testing.T) {
	db, err := Open("ramsql", "TestLoadUserAddresses")
	if err != nil {
		t.Fatalf("Open: error %v", err)
	}
	defer db.Close()

	s := db.Session()
	if s == nil {
		t.Fatalf("DB.Session: returns nil")
	}
	defer s.Close()

	var q string
	var res Result
	var rows *Rows

	q = "CREATE TABLE user (id INT PRIMARY KEY AUTOINCREMENT, name TEXT, age INT);"
	res, err = s.Exec(q)
	if err != nil {
		t.Fatalf("Session.Exec: %s error %v", q, err)
	}
	if res.LastInsertId != 0 {
		t.Errorf("Session.Exec: LastInsertId %d", res.LastInsertId)
	}
	if res.RowsAffected != 1 {
		t.Errorf("Session.Exec: RowsAffected %d", res.RowsAffected)
	}

	q = "INSERT INTO user (name, age) VALUES ('Akihabara', 32);"
	res, err = s.Exec(q)
	if err != nil {
		t.Fatalf("Session.Exec: %s error %v", q, err)
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
		res, err = s.Exec(q, name, age)
		if err != nil {
			t.Fatalf("Session.Exec: %s error %v", q, err)
		}
	}

	q = "SELECT id, name, age FROM user WHERE age = ?;"
	rows, err = s.Query(q, 32)
	if err != nil {
		t.Fatalf("Session.Query: %s error %v", q, err)
	}
	if rows == nil {
		t.Fatalf("Session.Query: %s returns nil", q)
	}

	scn := func(sc Scan) ([]interface{}, error) {
		id := s.Int64()
		name := s.String()
		age := s.Int64()
		err := sc(id, name, age)
		if err != nil {
			return nil, err
		}
		return []interface{}{id, name, age}, nil
	}

	row, err := rows.Fetch(scn)
	if err != nil {
		t.Fatalf("Rows.Fetch: error %v", err)
	}
	if row == nil {
		t.Fatalf("Rows.Fetch: returns nil")
	}
	if row.String("name") != "Akihabara" {
		t.Errorf("Row.String: expected Akihabara, got %s", row.String("name"))
	}

	row, err = rows.Fetch(scn)
	if err != nil {
		t.Fatalf("Rows.Fetch: error %v", err)
	}
	if row != nil {
		t.Fatalf("Rows.Fetch: returns %v", row)
	}
}
