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
		if err := b.buildNodeFile(node); err != nil {
			return err
		}
	}

	for _, object := range b.objects {
		if err := b.buildObjectFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *fieldBuilder) buildNodeFile(node *field.NodeTable) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     node,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Var().Id(node.Name).Op("=").Id("new" + node.Name).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Line()
	f.Add(b.fieldNew(node))

	f.Line()
	f.Type().Id(node.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").String())
			for _, fld := range node.Fields {
				if code := fld.CodeGen().FieldDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range node.GetFields() {
		if code := fld.CodeGen().FieldFunc(fieldCtx); code != nil {
			f.Line()
			f.Add(code)
		}
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *fieldBuilder) fieldNew(node *field.NodeTable) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     node,
	}

	return jen.Func().Id("new" + node.Name).
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").String()).
		Id(node.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(
				jen.Id(node.NameGoLower()).Types(def.TypeModel).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("key")] = jen.Id("key")
						for _, f := range node.Fields {
							if code := f.CodeGen().FieldInit(fieldCtx); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}

func (b *fieldBuilder) buildObjectFile(object *field.DatabaseObject) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     object,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Add(b.fieldNewObject(object))

	f.Line()
	f.Type().Id(object.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").String())
			for _, fld := range object.Fields {
				if code := fld.CodeGen().FieldDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range object.GetFields() {
		if code := fld.CodeGen().FieldFunc(fieldCtx); code != nil {
			f.Line()
			f.Add(code)
		}
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), object.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *fieldBuilder) fieldNewObject(object *field.DatabaseObject) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     object,
	}

	return jen.Func().Id("new" + object.Name).
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").String()).
		Id(object.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(
				jen.Id(object.NameGoLower()).Types(def.TypeModel).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("key")] = jen.Id("key")
						for i, f := range object.Fields {
							fCtx := fieldCtx
							if object.IsArrayIndexed {
								idx := i
								fCtx.ArrayIndex = &idx
							}
							if code := f.CodeGen().FieldInit(fCtx); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}
