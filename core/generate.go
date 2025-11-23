package core

import (
	"fmt"
	"github.com/go-surreal/som/core/codegen"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
	"github.com/go-surreal/som/core/util/gomod"
	"path"
	"path/filepath"
	"strings"
)

func Generate(inPath, outPath string, verbose, dry, check bool) error {
	absDir, err := filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("could not find absolute path: %v", err)
	}

	mod, err := gomod.FindGoMod(absDir)
	if err != nil {
		return fmt.Errorf("could not find go.mod: %v", err)
	}

	if check {
		info, err := mod.CheckGoVersion()
		if err != nil {
			return err
		}

		if verbose && info != "" {
			fmt.Println("ⓘ ", info)
		}
	}

	if check {
		info, err := mod.CheckSOMVersion(verbose)
		if err != nil {
			return err
		}

		if verbose && info != "" {
			fmt.Println("ⓘ ", info)
		}
	}

	info, err := mod.CheckDriverVersion()
	if err != nil {
		return err
	}

	if verbose && info != "" {
		fmt.Println("ⓘ ", info)
	}

	if err := mod.Save(); err != nil {
		return err
	}

	outPkg := path.Join(mod.Module(), strings.TrimPrefix(absDir, mod.Dir()))

	out := fs.New()

	err = codegen.BuildStatic(out, outPkg)
	if err != nil {
		return fmt.Errorf("could not generate code: %w", err)
	}

	if err := out.Flush(absDir); err != nil {
		return fmt.Errorf("could not write static files: %w", err)
	}

	source, err := parser.Parse(inPath, outPkg)
	if err != nil {
		return fmt.Errorf("could not parse source: %w", err)
	}

	err = codegen.Build(source, out, outPkg)
	if err != nil {
		return fmt.Errorf("could not generate code: %w", err)
	}

	if verbose {
		if err := out.Dry(absDir); err != nil {
			return err
		}
	}

	if dry {
		return nil
	}

	return out.Flush(absDir)
}
