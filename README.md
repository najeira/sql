# sql

A useful package of SQL.

This library wraps standard `database/sql` package
and provides some utilities.

## Usage

```go
import "github.com/najeira/sql"

func main() {
    // open
    db, err := sql.Open(sql.Config{
        User:       "sqltest",
        Passwd:     "testsql",
        ServerName: "localhost:3306",
        DBName:     "sqltest",
    })
    if err != nil {
        panic(err)
    }
    
    // select
    ctx := context.Background()
    var rows []*struct{
        ID int64 `db:"id"`
        Name string `db:"name"`
    }
    q := "SELECT id, name FROM users"
    res, err := db.Select(ctx, &rows, q)
    if err != nil {
        panic(err)
    }
}
```
