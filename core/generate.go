package core

import (
	"fmt"
	"github.com/go-surreal/som/core/codegen"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util"
	"github.com/go-surreal/som/core/util/fs"
	"path"
	"path/filepath"
	"strings"
)

func Generate(inPath, outPath string, verbose, dry, check bool) error {
	absDir, err := filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("could not find absolute path: %v", err)
	}

	mod, err := util.FindGoMod(absDir)
	if err != nil {
		return fmt.Errorf("could not find go.mod: %v", err)
	}

	if check {
		if info, err := mod.CheckGoVersion(); err != nil {
			return err
		} else if info != "" {
			fmt.Println("ⓘ ", info)
		}

		if info, err := mod.CheckSOMVersion(); err != nil {
			return err
		} else if info != "" {
			fmt.Println("ⓘ ", info)
		}

		if info, err := mod.CheckSDBCVersion(); err != nil {
			return err
		} else if info != "" {
			fmt.Println("ⓘ ", info)
		}
	}

	outPkg := path.Join(mod.Module(), strings.TrimPrefix(absDir, mod.Dir()))

	source, err := parser.Parse(inPath)
	if err != nil {
		return fmt.Errorf("could not parse source: %v", err)
	}

	out := fs.New()

	err = codegen.Build(source, out, outPkg)
	if err != nil {
		return fmt.Errorf("could not generate code: %v", err)
	}

	if verbose {
		if err := out.Dry(outPath); err != nil {
			return err
		}
	}

	if dry {
		return nil
	}

	return out.Flush(outPath)
}
