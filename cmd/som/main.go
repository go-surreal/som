package main

import (
	"context"
	"github.com/go-surreal/som/cmd/som/gen"
	cli "github.com/urfave/cli/v3"
	"log"
	"os"
	"runtime/debug"
)

func main() {
	ctx := context.Background()

	app := cli.Command{
		Name:  "somgen",
		Usage: "Generate SOM code for typesafe SurrealDB access",
		// ArgsUsage:      "<input_path> <output_path>",
		Description: "Tool for generating typesafe SurrealDB access layer from input models.",

		Commands: []*cli.Command{
			gen.Cmd(),
		},
		DefaultCommand: "gen",
		Suggest:        true,

		Authors: []any{
			"Marc Binz",
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

	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
