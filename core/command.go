package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/util/gomod"
	"github.com/urfave/cli/v3"
)

const (
	flagIn           = "in"
	flagOut          = "out"
	flagVerbose      = "verbose"
	flagDry          = "dry"
	flagNoCheck      = "no-check"
	flagNoCountIndex = "no-count-index"
	flagWire         = "wire"

	defaultOutputDir = "gen/som"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        "som",
		Usage:       "Generate code for typesafe SurrealDB access",
		Description: "Detects the closest go.mod and generates a typesafe access layer.\n\nOutput is written to 'gen/som' by default (relative to go.mod).\nIf --in is not specified, only static base files are generated.",
		Suggest:     true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagIn,
				Aliases: []string{"i"},
				Usage:   "Input model directory (relative to go.mod)",
			},
			&cli.StringFlag{
				Name:    flagOut,
				Aliases: []string{"o"},
				Usage:   "Output directory (relative to go.mod)",
				Value:   defaultOutputDir,
			},
			&cli.BoolFlag{
				Name:    flagVerbose,
				Aliases: []string{"v"},
			},
			&cli.BoolFlag{
				Name: flagDry,
			},
			&cli.BoolFlag{
				Name:  flagNoCheck,
				Usage: "Disable version checks for go and som",
			},
			&cli.BoolFlag{
				Name:  flagNoCountIndex,
				Usage: "Disable automatic COUNT index generation",
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
	if cmd.Args().Len() > 0 {
		return cli.Exit("unexpected positional arguments; use --in and --out flags instead", 1)
	}

	modDir, err := findModDir()
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if err := validateRelativePath(flagOut, cmd.String(flagOut)); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	outPath := filepath.Join(modDir, cmd.String(flagOut))

	init := !cmd.IsSet(flagIn)
	inPath := ""

	if !init {
		inDir := cmd.String(flagIn)
		if err := validateRelativePath(flagIn, inDir); err != nil {
			return cli.Exit(err.Error(), 1)
		}

		absInPath := filepath.Join(modDir, inDir)

		if _, err := os.Stat(absInPath); err != nil {
			return cli.Exit(fmt.Sprintf("input directory %q not found", inDir), 1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			return cli.Exit(fmt.Sprintf("could not get working directory: %v", err), 1)
		}

		inPath, err = filepath.Rel(cwd, absInPath)
		if err != nil {
			inPath = absInPath
		}
	}

	if err := Generate(inPath, outPath, init, cmd.Bool(flagVerbose), cmd.Bool(flagDry), !cmd.Bool(flagNoCheck), cmd.Bool(flagNoCountIndex), cmd.String(flagWire)); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}

func validateRelativePath(flagName, value string) error {
	if filepath.IsAbs(value) {
		return fmt.Errorf("--%s must be a relative path, got %q", flagName, value)
	}
	cleaned := filepath.Clean(value)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return fmt.Errorf("--%s must not escape the module root, got %q", flagName, value)
	}
	return nil
}

func findModDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %w", err)
	}

	mod, err := gomod.FindGoMod(cwd)
	if err != nil {
		return "", fmt.Errorf("could not find go.mod: %w", err)
	}

	return mod.Dir(), nil
}
