package gen

import (
	"context"
	"github.com/go-surreal/som/core"
	"github.com/urfave/cli/v3"
)

const (
	flagVerbose = "verbose"
	flagDry     = "dry"
	flagNoCheck = "no-check"
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
				Name:    flagVerbose,
				Aliases: []string{"v"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  flagDry,
				Value: false,
			},
			&cli.BoolFlag{
				Name:  flagNoCheck,
				Usage: "Disable version checks for go, som and sdbc",
			},
		},
		Action: generate,
	}
}

func generate(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 2 {
		return cli.Exit("Incorrect number of arguments", 1)
	}

	inPath := cmd.Args().Get(0)
	outPath := cmd.Args().Get(1)

	if err := core.Generate(inPath, outPath, cmd.Bool(flagVerbose), cmd.Bool(flagDry), !cmd.Bool(flagNoCheck)); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}
