package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
)

type filterBuilder struct {
	*baseBuilder
}

func newFilterBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *filterBuilder {
	return &filterBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *filterBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildFile(edge); err != nil {
			return err
		}
	}

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *filterBuilder) buildFile(elem field.Element) error {
	file := jen.NewFile(b.pkgName)

	file.PackageComment(string(embed.CodegenComment))

	if edge, ok := elem.(*field.EdgeTable); ok {
		b.buildEdge(file, edge)
	} else {
		b.buildOther(file, elem)
	}

	if err := file.Render(b.fs.Writer(path.Join(b.path(), elem.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *filterBuilder) buildOther(file *jen.File, elem field.Element) {
	pkgLib := b.subPkg(def.PkgLib)

	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	if _, ok := elem.(*field.NodeTable); ok {
		file.Line()
		file.Var().Id(elem.NameGo()).Op("=").
			Id("new" + elem.NameGo()).Types(b.SourceQual(elem.NameGo())).
			Call(jen.Qual(pkgLib, "NewKey").Types(b.SourceQual(elem.NameGo())).Call())
	}

	file.Line()
	file.Add(b.whereNew(elem))

	file.Type().Id(elem.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Qual(pkgLib, "Key").Types(def.TypeModel)) // TODO: name clash with Key field! -> go1.23: type key_[M any] = lib.Key[M]
			for _, f := range elem.GetFields() {
				if code := f.CodeGen().FilterDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range elem.GetFields() {
		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}

	// Generate extra filter code (wrapper types with Matches() for search-indexed fields)
	for _, fld := range elem.GetFields() {
		if code := fld.CodeGen().FilterExtra(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}

	// TODO: add record::exists filter function
	// https://github.com/surrealdb/surrealdb/pull/4602

	file.Line()
	file.Type().Id(elem.NameGoLower()+"Edges").
		Types(jen.Add(def.TypeModel).Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(def.TypeModel), // TODO: needed?
			jen.Qual(pkgLib, "Key").Types(def.TypeModel),
		)

	fieldCtx.Receiver = jen.Id(elem.NameGoLower() + "Edges").Types(def.TypeModel)

	for _, fld := range elem.GetFields() {
		isEdge := false

		if _, ok := fld.(*field.Edge); ok {
			isEdge = true
		}

		if slice, ok := fld.(*field.Slice); ok {
			if _, ok := slice.Element().(*field.Edge); ok {
				isEdge = true
			}
		}

		if !isEdge {
			continue
		}

		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}
}

func (b *filterBuilder) buildEdge(file *jen.File, edge *field.EdgeTable) {
	pkgLib := b.subPkg(def.PkgLib)

	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     edge,
	}

	file.Line()
	file.Var().Id(edge.NameGo()).Op("=").
		Id("new" + edge.NameGo()).Types(b.SourceQual(edge.NameGo())).
		Call(jen.Qual(pkgLib, "NewKey").Types(b.SourceQual(edge.NameGo())).Call())

	file.Line()
	file.Add(b.whereNew(edge))

	file.Line()
	file.Type().Id(edge.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Qual(pkgLib, "Key").Types(def.TypeModel))
			for _, f := range edge.GetFields() {
				if code := f.CodeGen().FilterDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range edge.GetFields() {
		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}

	// Generate extra filter code (wrapper types with Matches() for search-indexed fields)
	for _, fld := range edge.GetFields() {
		if code := fld.CodeGen().FilterExtra(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}

	file.Line()
	file.Type().Id(edge.NameGoLower()+"In").
		Types(jen.Add(def.TypeModel).Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(def.TypeModel),
			jen.Id("key").Qual(pkgLib, "Key").Types(def.TypeModel),
		)

	file.Line()
	file.Func().Id("new" + edge.NameGo() + "In").
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(def.TypeModel)).
		Id(edge.NameGoLower() + "In").Types(def.TypeModel).
		Block(
			jen.Return(
				jen.Id(edge.NameGoLower()+"In").Types(def.TypeModel).Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			),
		)

	file.Line()
	file.Func().
		Params(jen.Id("i").Id(edge.NameGoLower()+"In").Types(def.TypeModel)).Id(edge.Out.NameGo()).
		Params(
			jen.Id("filters").Op("...").Qual(pkgLib, "Filter").Types(b.SourceQual(edge.Out.Table().NameGo())),
		).
		Id(edge.Out.Table().NameGoLower()+"Edges").Types(def.TypeModel).
		Block(
			jen.Id("key").Op(":=").Qual(pkgLib, "EdgeIn").Call(
				jen.Id("i").Dot("key"),
				jen.Lit(edge.Out.NameDatabase()),
				jen.Id("filters"),
			),
			jen.Return(jen.Id(edge.Out.Table().NameGoLower()+"Edges").Types(def.TypeModel).
				Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			))

	file.Line()
	file.Type().Id(edge.NameGoLower()+"Out").
		Types(jen.Add(def.TypeModel).Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(def.TypeModel),
			jen.Id("key").Qual(pkgLib, "Key").Types(def.TypeModel),
		)

	file.Line()
	file.Func().Id("new" + edge.NameGo() + "Out").
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(def.TypeModel)).
		Id(edge.NameGoLower() + "Out").Types(def.TypeModel).
		Block(
			jen.Return(
				jen.Id(edge.NameGoLower()+"Out").Types(def.TypeModel).Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			),
		)

	file.Line()
	file.Func().
		Params(jen.Id("o").Id(edge.NameGoLower()+"Out").Types(def.TypeModel)).Id(edge.In.NameGo()).
		Params(
			jen.Id("filters").Op("...").Qual(pkgLib, "Filter").Types(b.SourceQual(edge.In.Table().NameGo())),
		).
		Id(edge.In.Table().NameGoLower()+"Edges").Types(def.TypeModel).
		Block(
			jen.Id("key").Op(":=").Qual(pkgLib, "EdgeOut").Call(
				jen.Id("o").Dot("key"),
				jen.Lit(edge.In.NameDatabase()),
				jen.Id("filters"),
			),
			jen.Return(jen.Id(edge.In.Table().NameGoLower()+"Edges").Types(def.TypeModel).
				Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			))
}

func (b *filterBuilder) whereNew(elem field.Element) jen.Code {
	pkgLib := b.subPkg(def.PkgLib)

	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	return jen.Func().Id("new" + elem.NameGo()).
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(def.TypeModel)).
		Id(elem.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(
				jen.Id(elem.NameGoLower()).Types(def.TypeModel).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("Key")] = jen.Id("key")
						for i, f := range elem.GetFields() {
							fCtx := fieldCtx
							if obj, ok := elem.(*field.DatabaseObject); ok && obj.IsArrayIndexed {
								idx := i
								fCtx.ArrayIndex = &idx
							}
							if code := f.CodeGen().FilterInit(fCtx); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}
