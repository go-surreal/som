package core

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/codegen"
	"github.com/go-surreal/som/core/parser"
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

	outPkg := path.Join(mod.Module(), strings.TrimPrefix(absDir, mod.Dir()))

	out := fs.New()

	var source *parser.Output
	usedFeatures := &parser.UsedFeatures{}

	if !init {
		// Parse first to determine which features are used.
		source, err = parser.Parse(inPath, outPkg,
			[]parser.TypeHandler{
				&parser.Node{},
				&parser.Union{},
				&parser.Edge{},
				&parser.View{},
				&parser.Sink{},
				&parser.Fragment{},
				&parser.ComplexIDStructHandler{},
				&parser.Enum{},
				&parser.EnumValue{},
				&parser.Struct{},
			},
			[]parser.FieldHandler{
				&parser.FieldEmail{},
				&parser.FieldSemVer{},
				&parser.FieldPassword{},
				&parser.FieldEnum{},
				&parser.FieldDuration{},
				&parser.FieldTime{},
				&parser.FieldMonth{},
				&parser.FieldWeekday{},
				&parser.FieldGeometry{},
				&parser.FieldURL{},
				&parser.FieldNode{},
				&parser.FieldUnion{},
				&parser.FieldEdge{},
				parser.InvalidView(),
				parser.InvalidSink(),
				parser.InvalidFragment(),
				&parser.FieldUUID{},
				&parser.FieldSlice{},
				&parser.FieldBool{},
				&parser.FieldByte{},
				&parser.FieldString{},
				&parser.FieldNumeric{},
				&parser.FieldStruct{},
			},
		)
		if err != nil {
			return fmt.Errorf("could not parse source: %w", err)
		}

		usedFeatures = source.UsedFeatures

		if err := checkLibVersions(mod, usedFeatures); err != nil {
			return err
		}
	}

	if err := mod.Save(); err != nil {
		return err
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

func checkLibVersions(mod *gomod.GoMod, features *parser.UsedFeatures) error {
	if features.UsesGoogleUUID {
		if err := mod.CheckLibVersion(gomod.PkgUUIDGoogle, gomod.MinUUIDGoogleVersion); err != nil {
			return err
		}
	}
	if features.UsesGofrsUUID {
		if err := mod.CheckLibVersion(gomod.PkgUUIDGofrs, gomod.MinUUIDGofrsVersion); err != nil {
			return err
		}
	}
	if features.UsesStdUUID {
		if err := mod.CheckStdUUIDSupport(); err != nil {
			return err
		}
	}
	if features.UsesOrbGeo {
		if err := mod.CheckLibVersion(gomod.PkgGeoOrb, gomod.MinGeoOrbVersion); err != nil {
			return err
		}
	}
	if features.UsesSimplefeaturesGeo {
		if err := mod.CheckLibVersion(gomod.PkgGeoSimplefeatures, gomod.MinGeoSimplefeaturesVersion); err != nil {
			return err
		}
	}
	if features.UsesGoGeomGeo {
		if err := mod.CheckLibVersion(gomod.PkgGeoGoGeom, gomod.MinGeoGoGeomVersion); err != nil {
			return err
		}
	}
	return nil
}
