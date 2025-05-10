package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
)

type convBuilder struct {
	*baseBuilder
}

func newConvBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *convBuilder {
	return &convBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *convBuilder) build() error {
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

func (b *convBuilder) buildFile(elem field.Element) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	typeName := elem.NameGoLower()
	if isNode || isEdge {
		typeName = elem.NameGo()
	}

	f.Line()
	f.Type().Id(typeName).StructFunc(func(g *jen.Group) {
		for _, f := range elem.GetFields() {
			if code := f.CodeGen().FieldDef(fieldCtx); code != nil {
				g.Add(code)
			}
		}
	})

	f.Line()
	f.Add(b.buildFrom(elem))

	f.Line()
	f.Add(b.buildTo(elem))

	if node, ok := elem.(*field.NodeTable); ok {
		f.Line()
		f.Type().Id(node.NameGoLower()+"Link").Struct(
			jen.Id(node.NameGo()),
			jen.Id("ID").Op("*").Qual(b.subPkg("sombase"), "ID"),
		)

		f.Line()
		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGoLower()+"Link")).
			Id("MarshalCBOR").Params().
			Params(jen.Index().Byte(), jen.Error()).
			Block(
				jen.If(jen.Id("f").Op("==").Nil()).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				),
				jen.Return(jen.Qual(def.PkgCBOR, "Marshal").Call(jen.Id("f").Dot("ID"))),
			)

		f.Line()
		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGoLower()+"Link")).
			Id("UnmarshalCBOR").Params(jen.Id("data").Index().Byte()).
			Error().
			Block(
				jen.If(
					jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("f").Dot("ID")),
					jen.Err().Op("==").Nil(),
				).Block(
					jen.Return(jen.Nil()),
				),

				jen.Type().Id("alias").Id(node.NameGoLower()+"Link"),
				jen.Var().Id("link").Id("alias"),

				jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("link")),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.Op("*").Id("f").Op("=").Id(node.NameGoLower()+"Link").Call(jen.Id("link")),
				),

				jen.Return(jen.Err()),
			)

		f.Line()
		f.Add(b.buildFromLink(node))

		f.Line()
		f.Add(b.buildFromLinkPtr(node))

		f.Line()
		f.Add(b.buildToLink(node))

		f.Line()
		f.Add(b.buildToLinkPtr(node))
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), elem.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *convBuilder) buildFrom(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	localName := elem.NameGoLower()
	methodPrefix := "from"

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	if isNode || isEdge {
		localName = elem.NameGo()
		methodPrefix = "From"
	}

	return jen.Add(
		// NO PTR
		jen.Func().
			Id(methodPrefix+elem.NameGo()).
			Params(jen.Id("data").Add(b.SourceQual(elem.NameGo()))).
			Id(localName).
			Block(
				jen.Return(jen.Id(localName).Values(jen.DictFunc(func(d jen.Dict) {
					for _, f := range elem.GetFields() {
						if code := f.CodeGen().ConvFrom(fieldCtx); code != nil {
							d[jen.Id(f.NameGo())] = code
						}
					}
				}))),
			),

		jen.Line(),

		// PTR
		jen.Func().
			Id(methodPrefix+elem.NameGo()+"Ptr").
			Params(jen.Id("data").Op("*").Add(b.SourceQual(elem.NameGo()))).
			Op("*").Id(localName).
			Block(
				jen.If(jen.Id("data").Op("==").Nil()).Block(
					jen.Return(jen.Nil()),
				),

				jen.Return(jen.Op("&").Id(localName).Values(jen.DictFunc(func(d jen.Dict) {
					for _, f := range elem.GetFields() {
						if code := f.CodeGen().ConvFrom(fieldCtx); code != nil {
							d[jen.Id(f.NameGo())] = code
						}
					}
				}))),
			),
	)
}

func (b *convBuilder) buildTo(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	localName := elem.NameGoLower()
	methodPrefix := "to"

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	if isNode || isEdge {
		localName = elem.NameGo()
		methodPrefix = "To"
	}

	ptr := jen.Empty()

	if isEdge {
		ptr = jen.Op("*")
	}

	return jen.Add(
		// NO PTR
		jen.Func().
			Id(methodPrefix+elem.NameGo()).
			Params(jen.Id("data").Add(ptr).Id(localName)).
			Add(b.SourceQual(elem.NameGo())).
			Block(
				jen.Return(jen.Add(b.SourceQual(elem.NameGo())).Values(jen.DictFunc(func(d jen.Dict) {
					for _, f := range elem.GetFields() {
						if code := f.CodeGen().ConvTo(fieldCtx); code != nil {
							if fieldCode := f.CodeGen().ConvToField(fieldCtx); fieldCode != nil {
								d[fieldCode] = code
								continue
							}

							d[jen.Id(f.NameGo())] = code
						}
					}

					if _, ok := elem.(*field.NodeTable); ok {
						d[jen.Id("Node")] = jen.Qual(b.subPkg("sombase"), "NewNode").Call(
							jen.Id("data").Dot("ID"),
						)
					}

					if _, ok := elem.(*field.EdgeTable); ok {
						d[jen.Id("Edge")] = jen.Qual(b.subPkg("sombase"), "NewEdge").Call(
							jen.Id("data").Dot("ID"),
						)
					}
				}))),
			),

		jen.Line(),

		// PTR
		jen.Func().
			Id(methodPrefix+elem.NameGo()+"Ptr").
			Params(jen.Id("data").Op("*").Id(localName)).
			Op("*").Add(b.SourceQual(elem.NameGo())).
			Block(
				jen.If(jen.Id("data").Op("==").Nil()).Block(
					jen.Return(jen.Nil()),
				),

				jen.Return(jen.Op("&").Add(b.SourceQual(elem.NameGo())).Values(jen.DictFunc(func(d jen.Dict) {
					for _, f := range elem.GetFields() {
						if code := f.CodeGen().ConvTo(fieldCtx); code != nil {
							if fieldCode := f.CodeGen().ConvToField(fieldCtx); fieldCode != nil {
								d[fieldCode] = code
								continue
							}

							d[jen.Id(f.NameGo())] = code
						}
					}

					if _, ok := elem.(*field.NodeTable); ok {
						d[jen.Id("Node")] = jen.Qual(b.subPkg("sombase"), "NewNode").Call(
							jen.Id("data").Dot("ID"),
						)
					}

					if _, ok := elem.(*field.EdgeTable); ok {
						d[jen.Id("Edge")] = jen.Qual(b.subPkg("sombase"), "NewEdge").Call(
							jen.Id("data").Dot("ID"),
						)
					}
				}))),
			),
	)
}

func (b *convBuilder) buildFromLink(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"Link").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower()+"Link")).
		Add(b.SourceQual(node.NameGo())).
		Block(
			jen.If(jen.Id("link").Op("==").Nil()).Block(
				jen.Return(jen.Add(b.SourceQual(node.NameGo())).Values()),
			),
			jen.Id("res").Op(":=").Id(node.NameGo()).Call(jen.Id("link").Dot(node.NameGo())),
			jen.Return(jen.Id("To"+node.NameGo()).Call(jen.Id("res"))),
		)
}

func (b *convBuilder) buildFromLinkPtr(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"LinkPtr").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower()+"Link")).
		Op("*").Add(b.SourceQual(node.NameGo())).
		Block(
			jen.If(jen.Id("link").Op("==").Nil()).Block(
				jen.Return(jen.Nil()),
			),
			jen.Id("res").Op(":=").Id(node.NameGo()).Call(jen.Id("link").Dot(node.NameGo())),
			jen.Id("out").Op(":=").Id("To"+node.NameGo()).Call(jen.Id("res")),
			jen.Return(jen.Id("&").Id("out")),
		)
}

func (b *convBuilder) buildToLink(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("to"+node.NameGo()+"Link").
		Params(jen.Id("node").Add(b.SourceQual(node.NameGo()))).
		Op("*").Id(node.NameGoLower()+"Link").
		Block(
			jen.If(jen.Id("node").Dot("ID").Call().Op("==").Nil()).Block(
				jen.Return(jen.Nil()),
			),
			jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
				jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(jen.Id("node")),
				jen.Id("ID").Op(":").Id("node").Dot("ID").Call(),
			),
			jen.Return(jen.Op("&").Id("link")),
		)
}

func (b *convBuilder) buildToLinkPtr(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("to"+node.NameGo()+"LinkPtr").
		Params(jen.Id("node").Op("*").Add(b.SourceQual(node.NameGo()))).
		Op("*").Id(node.NameGoLower()+"Link").
		Block(
			jen.
				If(
					jen.Id("node").Op("==").Nil().Op("||").
						Id("node").Dot("ID").Call().Op("==").Nil(),
				).
				Block(
					jen.Return(jen.Nil()),
				),
			jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
				jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(jen.Op("*").Id("node")),
				jen.Id("ID").Op(":").Id("node").Dot("ID").Call(),
			),
			jen.Return(jen.Op("&").Id("link")),
		)
}
