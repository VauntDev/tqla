# tqla

[![Vaunt Community](https://api.vaunt.dev/v1/github/entities/VauntDev/repositories/undefined/badges/community)](https://community.vaunt.dev/board/VauntDev/repository/tqla)

Tqla - (Template Query Language Agnostic)

Mas sql templating!!

## Install

`import "github.com/VauntDev/tqla"`

## Use

Tqla is a small lightweight text parser that wraps the golang `text/template` standard library. The primary purpose of tqla is to parse a text template and replace any variable with a placeholder. Variables that are replaced with placeholders are added to an args slice that can be passed to the standard db driver.

Similar libraries are exposed to sql injection, as they use the text/template library to do a simple text replacement. Tqla prevents sql injection by leveraging DB placeholders as described [here](https://go.dev/doc/database/sql-injection).

Currently, tqla does not attempt to do any sql validation, this may change in the future.

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

The test files also include examples of what can be done with tqla.

### Text Templating Basics

Tqla leverages golang text templating framework. Similar to the html/template package, tqla is a wrapper around text/templating that replaces values with db placeholders.

Below are a few common examples to get you started with Golang's text templating. More information can be found at [text/templating](https://pkg.go.dev/text/template).

#### Actions

Actions represent data evaluations, functions, conditions, and control loops. Actions are represented by `{{ <action> }}`.

#### Conditions

`if` conditionals are supported in templates.

Below is an example:

```go
{{ if .Value }}
  The value is {{ . Value }}
{{ end }}
```

or more complex uses:

```go
{{ if (gt .Value 0) }}
    The value, {{.Value}}, is great than 0
{{ else }}
    the value, {{.Value}}, is less than 0
{{ end }}
```

#### Loops

The `range` function is supported in templates.

below are a few examples

```go
// slice 
{{range .}}
{{.}}
{{end}}
```

or

```go
//map
{{range $key, $val := .}}
{{$key}} , {{$val}}
{{end}}
```

#### Functions

Go's text templating has several built-in functions outlined [here](https://pkg.go.dev/text/template#hdr-Functions). Tqla also supports additional functions defined by the user.

Here is an example of adding and calling a custom function within a template:

```go
funcs := template.FuncMap{
    "add": func(x int, y int) int {
        return x + y
    },
}

tqla, err := New(WithFuncMap(funcs))

s, a, err := tqla.Compile(`select * from table where v = {{ add . . }}`, 10)
if err != nil {
    t.Fatal(err)
}
```

Refer to the go doc for text/templating for a deeper dive into custom functions.

#### Sql Templates

A complete sql template may look something like the following:

```go
{{ define "nested_select" }}
    select count(*) as count_v
    from d.table as t
    join d.table_2 as t2 on t.c1 = t2.c2
    where t.c3 ilike {{ .FirstName }} and t.c4 = 'value'
{{end}}
{{ $maxLimit := 100 }}
select c1,
    c2,
    c3,
    COALESCE(c4,''),
    ({{ template "nested_select" . }})
from d.table as t
where c1 ilike ({{ .FirstName }})
{{ if .LastName }}
    and c2 = {{ .LastName }}
{{ end }}
{{ if .Page }}
    and (c3,c4) > ({{ .Page.Id }}, {{ .Page.Dttm }})
{{ end }}
{{ if gt .Page.Limit $maxLimit }}
    LIMIT {{ $maxLimit }}
{{ else }}
    LIMIT {{ .Page.Limit }}
{{ end }}`
```

## Contributor Guide

See our [Contributing.md](.github/CONTRIBUTING.md) file for guidance on how to contribute to tqla!

### Available Awards

| Name | Icon | Description |
| ---- | ---- | ----------- |
| Mas Tqla! | <img src="https://raw.githubusercontent.com/VauntDev/tqla/main/.vaunt/mas_tqla.png" width="150" /> | uno mas, dba..., one more, mas tqla! |
| Tqla Worm | <img src="https://raw.githubusercontent.com/VauntDev/tqla/main/.vaunt/tqla_worm.png" width="150" /> | not all bugs are bad! |
| Tqla Sunrise | <img src="https://raw.githubusercontent.com/VauntDev/tqla/main/.vaunt/tqla_sunrise.png" width="150" /> | a shining example of your contribution to tqla! |
| Tqla Honey Bee | <img src="https://raw.githubusercontent.com/VauntDev/tqla/main/.vaunt/tqla_honey_bee.png" width="150" /> | We bee-lieve in you! This achievement celebrates your first bug report, helping make Tqla sweeter than ever. |
| Mojito Tqla | <img src="https://raw.githubusercontent.com/VauntDev/tqla/main/.vaunt/mojito_tqla.png" width="150" /> | Your enhancement's a delight! Thanks for making Tqla even brighter than night. |

