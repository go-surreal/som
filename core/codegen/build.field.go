package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
)

type fieldBuilder struct {
	*baseBuilder
}

func newFieldBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *fieldBuilder {
	return &fieldBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *fieldBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *fieldBuilder) buildFile(node *field.NodeTable) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	modelType := b.SourceQual(node.NameGo())

	type fieldEntry struct {
		nameGo     string
		typeCode   jen.Code
		factoryFn  string
		dbName     string
	}

	var entries []fieldEntry

	for _, fld := range node.Fields {
		info := field.FieldDescriptorFor(fld)
		if info == nil {
			continue
		}

		entries = append(entries, fieldEntry{
			nameGo:    fld.NameGo(),
			typeCode:  info.TypeCode,
			factoryFn: info.FactoryName,
			dbName:    fld.NameDatabase(),
		})
	}

	if len(entries) == 0 {
		return nil
	}

	f := jen.NewFile(b.pkgName)
	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Var().Id(node.NameGo()).Op("=").
		StructFunc(func(g *jen.Group) {
			for _, e := range entries {
				g.Id(e.nameGo).Qual(pkgQuery, "Field").Types(modelType, e.typeCode)
			}
		}).
		Values(jen.DictFunc(func(d jen.Dict) {
			for _, e := range entries {
				var factoryCall *jen.Statement
				if e.factoryFn == "NewField" {
					factoryCall = jen.Qual(pkgQuery, e.factoryFn).
						Types(modelType, e.typeCode).
						Call(jen.Lit(e.dbName))
				} else {
					factoryCall = jen.Qual(pkgQuery, e.factoryFn).
						Types(modelType).
						Call(jen.Lit(e.dbName))
				}
				d[jen.Id(e.nameGo)] = factoryCall
			}
		}))

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
