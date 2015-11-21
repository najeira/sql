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
		t.Fatalf("DB.Session: nil")
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
		t.Fatalf("Session.Query: %s nil", q)
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
		t.Fatalf("Rows.Fetch: nil")
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

	row, err = rows.Fetch(scn)
	if err != nil {
		t.Fatalf("Rows.Fetch: error %v", err)
	}
	if row != nil {
		t.Fatalf("Rows.Fetch: %v", row)
	}
}

func TestSessionClose(t *testing.T) {
	db, err := Open("ramsql", "TestLoadUserAddresses")
	if err != nil {
		t.Fatalf("Open: error %v", err)
	}
	defer db.Close()

	s := db.Session()
	if s == nil {
		t.Fatalf("DB.Session: nil")
	}

	v := s.values

	if len(v.inuse) != 0 {
		t.Errorf("inuse: expected 0, got %d", len(v.inuse))
	}

	i := s.Int64()
	i.Valid = true
	i.Int64 = 123

	if len(v.inuse) != 1 {
		t.Errorf("inuse: expected 1, got %d", len(v.inuse))
	}

	err = s.Close()
	if err != nil {
		t.Errorf("Session.Close: %v", err)
	}

	if len(v.inuse) != 0 {
		t.Errorf("inuse: expected 0, got %d", len(v.inuse))
	}
}
