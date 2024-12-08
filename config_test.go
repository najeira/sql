package sql

import (
	"testing"
)

func TestConfigFormatDSN(t *testing.T) {
	var testDSNs = []struct {
		in  *Config
		out string
	}{
		{
			&Config{
				User:       "username",
				Passwd:     "password",
				ServerName: "localhost:3306",
				DBName:     "mydb",
			},
			"username:password@tcp(localhost:3306)/mydb?collation=utf8mb4_bin&interpolateParams=true",
		},
		{
			&Config{
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
