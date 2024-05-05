package gen

import (
	"github.com/go-surreal/som/core"
	"github.com/urfave/cli/v2"
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "gen",
		Aliases:     []string{"g"},
		Usage:       "Generate code for the database access based on input models",
		Description: "Takes the models from <input_path> and generates a typesafe access layer at <output_path>.",
		ArgsUsage:   "<input_path> <output_path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "nocheck",
				Usage: "Disable version checks for go, som and sdbc",
			},
		},
		Action: generate,
	}
}

func generate(ctx *cli.Context) error {
	if ctx.Args().Len() != 2 {
		return cli.Exit("Incorrect number of arguments", 1)
	}

	inPath := ctx.Args().Get(0)
	outPath := ctx.Args().Get(1)

	if err := core.Generate(inPath, outPath); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}
