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
