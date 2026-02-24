package codegen

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

const filenameInterfaces = "som.interfaces.go"

type build struct {
	input       *input
	fs          *fs.FS
	outPkg      string
	wirePackage string
}

func BuildStatic(fs *fs.FS, outPkg string, features *parser.UsedFeatures) error {
	tmpl := &embed.Template{
		GenerateOutPath: outPkg,
	}

	if features != nil {
		tmpl.UsesGoogleUUID = features.UsesGoogleUUID
		tmpl.UsesGofrsUUID = features.UsesGofrsUUID
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

func Build(source *parser.Output, fs *fs.FS, outPkg string, wirePackage string) error {
	in, err := newInput(source, outPkg)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:       in,
		fs:          fs,
		outPkg:      outPkg,
		wirePackage: wirePackage,
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
		if err := b.buildBaseFile(node); err != nil {
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

func (b *build) buildInterfaceFile() error {
	f := jen.NewFile(def.PkgRepo)

	f.PackageComment(string(embed.CodegenComment))

	f.Type().Id("Client").InterfaceFunc(func(g *jen.Group) {
		for _, node := range b.input.nodes {
			g.Id(node.NameGo() + "Repo").Call().Id(node.NameGo() + "Repo")
		}

		g.Id("ApplySchema").Call(jen.Id("ctx").Qual("context", "Context")).Error()
		g.Id("Close").Call()
	})

	// Generate ClientImpl with per-node cached repo fields.
	f.Line().Type().Id("ClientImpl").StructFunc(func(g *jen.Group) {
		g.Id("db").Id("Database")
		g.Id("mu").Qual("sync", "Mutex")
		for _, node := range b.input.nodes {
			g.Id(node.NameGoLower() + "Repo").Op("*").Id(node.NameGoLower())
		}
	})

	if err := f.Render(b.fs.Writer(filepath.Join(def.PkgRepo, filenameInterfaces))); err != nil {
		return err
	}

	return nil
}

func (b *build) buildBaseFile(node *field.NodeTable) error {
	return b.buildBaseFileFromTemplate(node)
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

func (b *build) newFieldBuilder() builder {
	return newFieldBuilder(b.input, b.fs, b.basePkg(), def.PkgField)
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) relativePkgPath(pkg ...string) string {
	return path.Join(append([]string{b.basePkg()}, pkg...)...)
}

func idFuncName(node *field.NodeTable) string {
	switch node.Source.IDType {
	case parser.IDTypeUUID:
		return "newUUID"
	case parser.IDTypeRand:
		return "newID"
	default:
		return "newULID" // ULID is the default ID type (used by the Node alias)
	}
}
