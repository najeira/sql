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
	if db == nil {
		t.Fatalf("Open: nil")
	}
	db.Close()
}
