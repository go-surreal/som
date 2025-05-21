package main

import (
	"context"
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

	theParser, err := parser.Parse(context.Background(), "../internal/tests/basic/model")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----------------------------")

	fmt.Println(theParser.String())
}
