# tqla

[![Vaunt Community](https://api.vaunt.dev/v1/github/entities/VauntDev/repositories/undefined/badges/community)](https://community.vaunt.dev/board/VauntDev/repository/undefined)

mas sql templating!

## Install

`import "github.com/VauntDev/tqla"`

## Use

Tqla is a small leight weight text parser that wraps the golang `text/template` standard library. The primary purpose of tqla is to parse a text template and replace any variable with a placeholder. Variables that are replaced with placeholders are added to an args slice that can be passed to standard db driver.

Currently, tqla does no sql validation, this may change in the future.

Below is a simple example on how it can be used:

```go
package main

import (
    "database/sql"
    "log"

    "github.com/VauntDev/tqla"
)

type example struct {
    Value string
}

func main() {

    db, err := sql.Open("mysql",
        "user:password@tcp(127.0.0.1:3306)/example")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    example := &example{Value: "example"}

    t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
    if err != nil {
        log.Fatal(err)
    }

    statement, args, err := t.Compile(`select *
    from table where Value = {{ .Value }}
    `, example)
    if err != nil {
        log.Fatal(err)
    }

    rows, err := db.Query(statement, args)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
}
```

### Simple Statements

holder
