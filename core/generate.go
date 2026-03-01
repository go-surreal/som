package core

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/codegen"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/fieldtype"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/go-surreal/som/core/util/fs"
	"github.com/go-surreal/som/core/util/gomod"
)

func Generate(inPath, outPath string, init, verbose, dry, check, noCountIndex bool, wireOverride string) error {
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

	var wirePackage string

	switch wireOverride {
	case "no":
		wirePackage = ""
	case "google":
		wirePackage = "github.com/google/wire"
	case "goforj":
		wirePackage = "github.com/goforj/wire"
	default:
		wirePackage = mod.WirePackage()
	}

	if err := mod.Save(); err != nil {
		return err
	}

	outPkg := path.Join(mod.Module(), strings.TrimPrefix(absDir, mod.Dir()))

	out := fs.New()

	var source *parser.Output
	usedFeatures := &parser.UsedFeatures{}

	if !init {
		// Parse first to determine which features are used.
		source, err = parser.Parse(inPath, outPkg,
			[]parser.TypeHandler{
				&structtype.NodeHandler{},
				&structtype.EdgeHandler{},
				&structtype.ComplexIDStructHandler{},
				&structtype.EnumHandler{},
				&structtype.EnumValueHandler{},
				&structtype.StructHandler{},
			},
			[]parser.FieldHandler{
				&fieldtype.EmailHandler{},
				&fieldtype.PasswordHandler{},
				&fieldtype.EnumHandler{},
				&fieldtype.DurationHandler{},
				&fieldtype.TimeHandler{},
				&fieldtype.MonthHandler{},
				&fieldtype.WeekdayHandler{},
				&fieldtype.GeometryHandler{},
				&fieldtype.URLHandler{},
				&fieldtype.NodeRefHandler{},
				&fieldtype.EdgeRefHandler{},
				&fieldtype.UUIDHandler{},
				&fieldtype.SliceHandler{},
				&fieldtype.BoolHandler{},
				&fieldtype.ByteHandler{},
				&fieldtype.StringHandler{},
				&fieldtype.NumericHandler{},
				&fieldtype.StructHandler{},
			},
		)
		if err != nil {
			return fmt.Errorf("could not parse source: %w", err)
		}

		usedFeatures = source.UsedFeatures
	}

	// Generate static files with feature flags from parsing.
	err = codegen.BuildStatic(out, outPkg, usedFeatures)
	if err != nil {
		return fmt.Errorf("could not generate code: %w", err)
	}

	if err := out.Flush(absDir); err != nil {
		return fmt.Errorf("could not write static files: %w", err)
	}

	if init {
		return nil
	}

	err = codegen.Build(source, out, outPkg, wirePackage, noCountIndex)
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
