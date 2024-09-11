package main

import (
	"fmt"
	"github.com/go-surreal/som/exp/parser"
	"log"
)

func main() {
	theParser := parser.NewParser()

	if err := theParser.Parse("./model"); err != nil {
		log.Fatal(err)
	}

	log.Println("-----------------------------")

	fmt.Println(theParser.String())
}
