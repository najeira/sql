# sql

A useful package of SQL.

This library wraps standard `database/sql` package
and provides some utilities.

## Usage

```go
import "github.com/najeira/sql"

func main() {
    // open
    db, err := sql.Open("mysql", "your_dsn")
    if err != nil {
        panic(err)
    }
    
    // query
    res, err := db.Query("SELECT id, name FROM person")
    if err != nil {
        panic(err)
    }
    defer res.Close()
    
    // fetch
    rows := make([]sql.Row)
    for res.Next() {
        var id sql.NullInt64
        var name sql.NullString
        row, err := res.Scan(&id, &name)
        if err != nil {
            return nil, err
        }
        rows = append(rows, row)
    }
    if err := res.Err(); err != nil {
        panic(err)
    }
    
    // use rows
    for _, row := range rows {
        // sql.Row has methods to get value by name
        fmt.Printf("id=%d, name=%s",
            row.Int64("id"),
            row.String("name"))
    }
}
```
