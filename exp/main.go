package main

import (
	"github.com/marcbinz/som/exp/parser"
	"log"
)

func main() {
	err := parser.Parse("./model")
	if err != nil {
		log.Fatal(err)
	}
}
