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
    _ "github.com/mattn/go-sqlite3"
)

type todo struct {
    Id          int
    Title       string
    Description string
    Completed   bool
}

const db = "example.db"
const todoSchema = `create table if not exists todos (
    id integer primary key,
    title text not null,
    description text not null,
    completed boolean default 0
);
`

func main() {

    log.Println("connecting to db... ")

    db, err := sql.Open("sqlite3", db)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    log.Println("creating table is it does not exist... ")

    if _, err := db.Exec(todoSchema); err != nil {
        log.Fatal(err)
    }

    todos := []*todo{
        {Id: 1, Title: "todo 1", Description: "first todo", Completed: false},
        {Id: 2, Title: "todo 2", Description: "second todo", Completed: false},
        {Id: 3, Title: "todo 3", Description: "third todo", Completed: false},
        {Id: 4, Title: "todo 4", Description: "fourth todo", Completed: false},
        {Id: 5, Title: "todo 5", Description: "fith todo", Completed: false},
    }

    t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
    if err != nil {
        log.Fatal(err)
    }

    log.Println("adding todos...")

    insertStmt, insertArgs, err := t.Compile(`
    {{ $len := 4 -}}
    INSERT INTO 'todos' ('id', 'title', 'description', 'completed')
    VALUES {{ range $i, $v := . }}
        ( {{$v.Id}}, {{$v.Title}}, {{$v.Description}}, {{ $v.Completed }} ){{if lt $i $len}},{{else}};{{end -}}
    {{end}}
    `, todos)
    if err != nil {
        log.Fatal(err)
    }

    if _, err := db.Exec(insertStmt, insertArgs...); err != nil {
        log.Fatal(err)
    }

    log.Println("looking up todo...")
    selectStmt, selectArgs, err := t.Compile(`select * from todos where id={{ . }}`, 5)
    if err != nil {
        log.Fatal(err)
    }

    todo := &todo{}
    row := db.QueryRow(selectStmt, selectArgs...)
    if err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
        log.Fatal(err)
    }

    log.Println("found: ", todo)
}
```

### Text/Template
