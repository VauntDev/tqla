package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/VauntDev/tqla"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "root"
	dbname   = "mysql-test-db"
)

type todo struct {
	Id          int
	Title       string
	Description string
	Completed   bool
}

const todoSchema = `CREATE TABLE IF NOT EXISTS todos (
	id INT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	completed BOOLEAN DEFAULT FALSE
);`

func main() {

	log.Println("connecting to db...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("creating table if it does not exist...")
	if _, err := db.Exec(todoSchema); err != nil {
		log.Fatal(err)
	}

	todos := []*todo{
		{Id: 1, Title: "todo 1", Description: "first todo", Completed: false},
		{Id: 2, Title: "todo 2", Description: "second todo", Completed: false},
		{Id: 3, Title: "todo 3", Description: "third todo", Completed: false},
		{Id: 4, Title: "todo 4", Description: "fourth todo", Completed: false},
		{Id: 5, Title: "todo 5", Description: "fifth todo", Completed: false},
	}

	t, err := tqla.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("adding todos...")
	insertStmt, insertArgs, err := t.Compile(`
	{{ $len := 4 -}}
	INSERT INTO todos (id, title, description, completed) 
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
	selectStmt, selectArgs, err := t.Compile(`SELECT * FROM todos WHERE id = {{ . }}`, 5)
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
