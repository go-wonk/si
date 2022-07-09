# si(storage interface)
`si` is a collection of wrappers that aims to ease reading/writing data from/to repositories. It is mostly a client side library and the following repositories or communication protocols will be supported from standard or non-standard packages.

- file
- tcp
- sql
- http
- websocket ([gorilla](https://github.com/gorilla/websocket))
- kafka ([sarama](https://github.com/Shopify/sarama))
- elasticsearch ([go-elasticsearch](https://github.com/elastic/go-elasticsearch))
- ftp ([jlaffaye](https://github.com/jlaffaye/ftp))

## Installation
```bash
go get -u github.com/go-wonk/si
```

## Quick Start
1. sql
```go
connStr := "host=testpghost port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
driver := "postgres"
db, _ := sql.Open(driver, connStr)
sqldb := sisql.NewSqlDB(db).WithTagKey("si")

type BoolTest struct {
    Nil      string `json:"nil" si:"nil"`
    True_1   bool   `json:"true_1" si:"true_1"`
    True_2   bool   `json:"true_2" si:"true_2"`
    False_1  bool   `json:"false_1" si:"false_1"`
    False_2  bool   `json:"false_2" si:"false_2"`
    Ignore_3 string `si:"-"`
}
query := `
    select null as nil,
        null as true_1, '1' as true_2, 
        0 as false_1, '0' as false_2
    union all
    select null as nil,
        1 as true_1, '1' as true_2,
        0 as false_1, '0' as false_2
`

m := []BoolTest{}
_, err := sqldb.QueryStructs(query, &m)

```
## Versions
### v0.1.1
- `siwrap` package has been renamed to `sisql`.