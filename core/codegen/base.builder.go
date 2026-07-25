package codegen

import (
	"fmt"
	"path"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

type build struct {
	input        *input
	fs           *fs.FS
	outPkg       string
	wirePackage  string
	noCountIndex bool
}

func BuildStatic(fs *fs.FS, outPkg string, features *parser.UsedFeatures) error {
	tmpl := &embed.Template{
		GenerateOutPath: outPkg,
	}

	if features != nil {
		tmpl.UsesGoogleUUID = features.UsesGoogleUUID
		tmpl.UsesGofrsUUID = features.UsesGofrsUUID
		tmpl.UsesOrbGeo = features.UsesOrbGeo
		tmpl.UsesSimplefeaturesGeo = features.UsesSimplefeaturesGeo
		tmpl.UsesGoGeomGeo = features.UsesGoGeomGeo
	}

	files, err := embed.Read(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		fs.Write(file.Path, file.Content)
	}

	return nil
}

func Build(source *parser.Output, fs *fs.FS, outPkg string, wirePackage string, noCountIndex bool) error {
	in, err := newInput(source, outPkg)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:        in,
		fs:           fs,
		outPkg:       outPkg,
		wirePackage:  wirePackage,
		noCountIndex: noCountIndex,
	}

	return builder.build()
}

func (b *build) build() error {
	if err := b.buildInterfaceFile(); err != nil {
		return err
	}

	if err := b.buildSchemaFile(); err != nil {
		return err
	}

	for _, node := range b.input.nodes {
		if err := b.buildNodeRepoFile(node); err != nil {
			return err
		}
	}

	for _, view := range b.input.views {
		if err := b.buildViewRepoFile(view); err != nil {
			return err
		}
	}

	for _, sink := range b.input.sinks {
		if err := b.buildSinkRepoFile(sink); err != nil {
			return err
		}
	}

	builders := []builder{
		b.newQueryBuilder(),
		b.newFilterBuilder(),
		b.newSortBuilder(),
		b.newFetchBuilder(),
		b.newConvBuilder(),
		b.newRelateBuilder(),
		b.newIndexBuilder(),
		b.newFieldBuilder(),
	}

	for _, builder := range builders {
		if err := builder.build(); err != nil {
			return err
		}
	}

	if b.wirePackage != "" {
		if err := b.buildWireFile(); err != nil {
			return err
		}
	}

	return nil
}

func (b *build) newQueryBuilder() builder {
	return newQueryBuilder(b.input, b.fs, b.basePkg(), def.PkgQuery)
}

func (b *build) newFilterBuilder() builder {
	return newFilterBuilder(b.input, b.fs, b.basePkg(), def.PkgFilter)
}

func (b *build) newSortBuilder() builder {
	return newSortBuilder(b.input, b.fs, b.basePkg(), def.PkgSort)
}

func (b *build) newFetchBuilder() builder {
	return newFetchBuilder(b.input, b.fs, b.basePkg(), def.PkgFetch)
}

func (b *build) newConvBuilder() builder {
	return newConvBuilder(b.input, b.fs, b.basePkg(), def.PkgConv)
}

func (b *build) newRelateBuilder() builder {
	return newRelateBuilder(b.input, b.fs, b.basePkg(), def.PkgRelate)
}

func (b *build) newIndexBuilder() builder {
	return newIndexBuilder(b.input, b.fs, b.basePkg(), def.PkgIndex, b.noCountIndex, b.input.define)
}

func (b *build) newFieldBuilder() builder {
	return newFieldBuilder(b.input, b.fs, b.basePkg(), def.PkgField)
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) relativePkgPath(pkg ...string) string {
	return path.Join(append([]string{b.basePkg()}, pkg...)...)
}
