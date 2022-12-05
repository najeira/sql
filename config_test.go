package sql_test

import (
	"testing"

	"github.com/najeira/sql"
)

func TestConfigFormatDSN(t *testing.T) {
	var testDSNs = []struct {
		in  *sql.Config
		out string
	}{
		{
			&sql.Config{
				User:       "username",
				Passwd:     "password",
				ServerName: "localhost:3306",
				DBName:     "mydb",
			},
			"username:password@tcp(localhost:3306)/mydb?collation=utf8mb4_bin&interpolateParams=true",
		},
		{
			&sql.Config{
				User:       "username",
				Passwd:     "password",
				ServerName: "/tmp/mysql",
				DBName:     "mydb",
			},
			"username:password@unix(/tmp/mysql)/mydb?collation=utf8mb4_bin&interpolateParams=true",
		},
	}
	for _, d := range testDSNs {
		got := d.in.FormatDSN()
		if got != d.out {
			t.Errorf(got)
		}
	}
}
