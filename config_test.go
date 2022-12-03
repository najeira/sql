package sql_test

import (
	"testing"

	"github.com/najeira/sql"
)

var testDSNs = []struct {
	in  *sql.Config
	out string
}{{
	&sql.Config{User: "username", Passwd: "password", ServerName: "localhost:3306", DBName: "mydb"},
	"username:password@tcp(localhost:3306)/mydb?collation=utf8mb4_bin&interpolateParams=true",
}}

func TestConfigFormatDSN(t *testing.T) {
	for _, d := range testDSNs {
		got := d.in.FormatDSN()
		if got != d.out {
			t.Errorf(got)
		}
	}
}
