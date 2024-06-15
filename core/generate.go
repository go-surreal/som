package core

import (
	"fmt"
	"github.com/go-surreal/som/core/codegen"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Generate(inPath, outPath string) error {
	absDir, err := filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("could not find absolute path: %v", err)
	}

	mod, err := util.FindGoMod(absDir)
	if err != nil {
		return fmt.Errorf("could not find go.mod: %v", err)
	}

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

	source, err := parser.Parse(inPath)
	if err != nil {
		return fmt.Errorf("could not parse source: %v", err)
	}

	if err := os.RemoveAll(outPath); err != nil {
		return err
	}

	diff := strings.TrimPrefix(absDir, mod.Dir())
	outPkg := path.Join(mod.Module(), diff)

	err = codegen.Build(source, outPath, outPkg)
	if err != nil {
		return fmt.Errorf("could not generate code: %v", err)
	}

	return nil
}
