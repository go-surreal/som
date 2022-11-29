package main

import (
	"github.com/marcbinz/som/cmd/somgen/sub"
	"github.com/marcbinz/som/core"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:   "somgen",
		Action: generate,
		Commands: []*cli.Command{
			sub.Surreal(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generate(ctx *cli.Context) error {
	inPath := ctx.Args().Get(0)
	outPath := ctx.Args().Get(1)

	return core.Generate(inPath, outPath)
}
