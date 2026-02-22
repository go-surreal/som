package core

import (
	"context"

	"github.com/urfave/cli/v3"
)

const (
	flagInit    = "init"
	flagVerbose = "verbose"
	flagDry     = "dry"
	flagNoCheck = "no-check"
	flagWire    = "wire"
)

func Gen() *cli.Command {
	return &cli.Command{
		Name:        "gen",
		Aliases:     []string{"g"},
		Usage:       "Generate code for the database access based on input models",
		Description: "Takes the models from <input_path> and generates a typesafe access layer at <output_path>.",
		ArgsUsage:   "<input_path> <output_path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  flagInit,
				Usage: `Initialize a new project with only the basic som files.`,
			},
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
				Usage: "Disable version checks for go and som",
			},
			&cli.StringFlag{
				Name:  flagWire,
				Usage: `Override wire generation: "no" to disable, "google" for github.com/google/wire, "goforj" for github.com/goforj/wire (default: auto-detect from go.mod)`,
			},
		},
		Action: generate,
	}
}

func generate(_ context.Context, cmd *cli.Command) error {
	init := cmd.Bool(flagInit)

	if (init && cmd.Args().Len() != 1) ||
		(!init && cmd.Args().Len() != 2) {
		return cli.Exit("Incorrect number of arguments", 1)
	}

	inPath := cmd.Args().Get(0)
	outPath := cmd.Args().Get(1)

	if init {
		inPath = "<not-used>"
		outPath = cmd.Args().Get(0)
	}

	if err := Generate(inPath, outPath, init, cmd.Bool(flagVerbose), cmd.Bool(flagDry), !cmd.Bool(flagNoCheck), cmd.String(flagWire)); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}
