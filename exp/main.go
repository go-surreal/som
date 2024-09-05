package main

import (
	"github.com/go-surreal/som/exp/parser"
	"log"
)

func main() {
	err := parser.NewParser().Parse("./model")
	if err != nil {
		log.Fatal(err)
	}
}
