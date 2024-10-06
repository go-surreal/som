package main

import (
	"fmt"
	"github.com/go-surreal/som/exp/parser"
	"log"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		),
	)

	theParser := parser.NewParser()

	if err := theParser.Parse("./model"); err != nil {
		log.Fatal(err)
	}

	log.Println("-----------------------------")

	fmt.Println(theParser.String())
}
