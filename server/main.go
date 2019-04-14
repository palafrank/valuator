package main

import (
	"log"

	"github.com/palafrank/valuator"
)

func main() {
	s, err := newServer("./db/", valuator.FileDatabaseType)
	if err == nil {
		log.Fatal(s.createAndRunServer("localhost", "8080"))
	}
	log.Fatal(err.Error())
}
