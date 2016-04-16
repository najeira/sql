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
    rows, err := res.Fetch(func(s sql.Scan){
        var id sql.NullInt64
        var name sql.NullString
        
        err := s.Scan(&id, &name)
        if err != nil {
            return nil, err
        }
        
        return []interface{&id, &name}, nil
    })
    if err != nil {
        panic(err)
    }
    
    // rows is []sql.Row
    for _, row := range rows {
        // sql.Row has methods to get value by name
        fmt.Printf("id=%d, name=%s",
            row.Int64("id"),
            row.String("name"))
    }
}
```
