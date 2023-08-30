package main

import (
	"log"

	"github.com/VauntDev/tqla"
)

type tester struct {
	Provider string
}

func main() {
	test := &tester{Provider: "github"}
	t, err := tqla.New(tqla.WithPlaceHolder(tqla.Dollar))
	if err != nil {
		log.Fatal(err)
	}
	s, a, err := t.Compile(`select * from vaunt.entities
	where provider = {{ .Provider }} and {{ .Provider }} = any(contributions.Labels)
	`, test)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s, a)
}
