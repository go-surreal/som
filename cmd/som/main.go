package main

import (
	"github.com/go-surreal/som/cmd/som/gen"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime/debug"
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

		Authors: []*cli.Author{
			{
				Name: "Marc Binz",
			},
		},
		Copyright: "github.com/go-surreal/som",

		ExtraInfo: func() map[string]string {
			info, ok := debug.ReadBuildInfo()
			if !ok {
				return nil
			}

			return map[string]string{
				"GoVersion": info.GoVersion,
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
