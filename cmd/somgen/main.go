package main

import (
	"github.com/go-surreal/som/buildtime"
	"github.com/go-surreal/som/cmd/somgen/gen"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:  "somgen",
		Usage: "Generate SOM code for typesafe SurrealDB access",
		// ArgsUsage:      "<input_path> <output_path>",
		Description: "Tool for generating typesafe SurrealDB access layer from input models.",

		Commands: []*cli.Command{
			gen.Cmd(),
		},
		DefaultCommand: "gen",
		Suggest:        true,

		Version:  buildtime.Version(),
		Compiled: buildtime.CompiledAt(),

		Authors: []*cli.Author{
			{
				Name: "Marc Binz",
			},
		},
		Copyright: "github.com/go-surreal/som",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
