package main

import (
	"github.com/marcbinz/sdb/core"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:   "sdbgen",
		Action: generate,
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
