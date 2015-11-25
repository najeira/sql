package sql

import (
	"testing"
)

func BenchmarkRowsWithPool(b *testing.B) {
	columns := []string{"foo", "bar"}

	disableRowsPool = false

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := getRowsForSqlRowsAndColumns(nil, columns)
		rows.Close()
	}
}

func BenchmarkRowsWithoutPool(b *testing.B) {
	columns := []string{"foo", "bar"}

	disableRowsPool = true

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows := getRowsForSqlRowsAndColumns(nil, columns)
		rows.Close()
	}
}
