package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
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
	f.Type().Id(typeName).Struct(
		jen.Add(b.SourceQual(elem.NameGo())),
	)

	f.Line()
	f.Add(b.buildMarshalCBOR(elem, typeName, fieldCtx, isNode, isEdge))

	f.Line()
	f.Add(b.buildUnmarshalCBOR(elem, typeName, fieldCtx, isNode, isEdge))

	f.Line()
	f.Add(b.buildFrom(elem))

	f.Line()
	f.Add(b.buildTo(elem))

	if node, ok := elem.(*field.NodeTable); ok {
		f.Line()
		f.Type().Id(node.NameGoLower()+"Link").Struct(
			jen.Id(node.NameGo()),
			jen.Id("ID").Op("*").Qual(b.subPkg(""), "ID"),
		)

		f.Line()
		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGoLower()+"Link")).
			Id("MarshalCBOR").Params().
			Params(jen.Index().Byte(), jen.Error()).
			Block(
				jen.If(jen.Id("f").Op("==").Nil()).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				),
				jen.Return(jen.Qual(path.Join(b.basePkg, "internal/cbor"), "Marshal").Call(jen.Id("f").Dot("ID"))),
			)

		f.Line()
		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGoLower()+"Link")).
			Id("UnmarshalCBOR").Params(jen.Id("data").Index().Byte()).
			Error().
			Block(
				jen.If(
					jen.Err().Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("f").Dot("ID")),
					jen.Err().Op("==").Nil(),
				).Block(
					jen.Return(jen.Nil()),
				),

				jen.Type().Id("alias").Id(node.NameGoLower()+"Link"),
				jen.Var().Id("link").Id("alias"),

				jen.Err().Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("link")),
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
	localName := elem.NameGoLower()
	methodPrefix := "from"

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	if isNode || isEdge {
		localName = elem.NameGo()
		methodPrefix = "From"
	}

	return jen.Add(
		// NO PTR - shallow wrapper: just embed
		jen.Func().
			Id(methodPrefix+elem.NameGo()).
			Params(jen.Id("data").Add(b.SourceQual(elem.NameGo()))).
			Id(localName).
			Block(
				jen.Return(jen.Id(localName).Values(jen.Dict{
					jen.Id(elem.NameGo()): jen.Id("data"), // ONE field copy
				})),
			),

		jen.Line(),

		// PTR - shallow wrapper: just embed
		jen.Func().
			Id(methodPrefix+elem.NameGo()+"Ptr").
			Params(jen.Id("data").Op("*").Add(b.SourceQual(elem.NameGo()))).
			Op("*").Id(localName).
			Block(
				jen.If(jen.Id("data").Op("==").Nil()).Block(
					jen.Return(jen.Nil()),
				),

				jen.Return(jen.Op("&").Id(localName).Values(jen.Dict{
					jen.Id(elem.NameGo()): jen.Op("*").Id("data"), // ONE field copy
				})),
			),
	)
}

func (b *convBuilder) buildMarshalCBOR(elem field.Element, typeName string, ctx field.Context, isNode, isEdge bool) jen.Code {
	return jen.Func().
		Params(jen.Id("c").Op("*").Id(typeName)).
		Id("MarshalCBOR").Params().
		Params(jen.Index().Byte(), jen.Error()).
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("c").Op("==").Nil()).Block(
				jen.Return(jen.Qual(path.Join(b.basePkg, "internal/cbor"), "Marshal").Call(jen.Nil())),
			)

			// Count fields for pre-sized map allocation
			fieldCount := 0
			if isNode || isEdge {
				fieldCount++ // ID field
			}
			for _, f := range elem.GetFields() {
				if f.NameDatabase() != "id" {
					fieldCount++
				}
			}

			g.Id("data").Op(":=").Make(jen.Map(jen.String()).Any(), jen.Lit(fieldCount))

			// Marshal ID field for nodes and edges
			if isNode || isEdge {
				tableName := elem.NameDatabase()
				g.Line()
				g.Comment("Embedded som.Node/Edge ID field")
				g.If(jen.Id("c").Dot("ID").Call().Op("!=").Lit("")).Block(
					jen.Id("data").Index(jen.Lit("id")).Op("=").Qual("github.com/surrealdb/surrealdb.go/pkg/models", "NewRecordID").Call(
						jen.Lit(tableName), jen.Id("c").Dot("ID").Call(),
					),
				)
			}

			// Marshal all fields
			g.Line()
			for _, f := range elem.GetFields() {
				// Skip ID field (handled specially for nodes/edges)
				if f.NameDatabase() == "id" {
					continue
				}

				// Generate marshal code for this field using field's CodeGen method
				if code := f.CodeGen().CBORMarshal(ctx); code != nil {
					g.Add(code)
				}
			}

			g.Line()
			g.Return(jen.Qual(path.Join(b.basePkg, "internal/cbor"), "Marshal").Call(jen.Id("data")))
		})
}

func (b *convBuilder) buildUnmarshalCBOR(elem field.Element, typeName string, ctx field.Context, isNode, isEdge bool) jen.Code {
	return jen.Func().
		Params(jen.Id("c").Op("*").Id(typeName)).
		Id("UnmarshalCBOR").Params(jen.Id("data").Index().Byte()).
		Error().
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("rawMap").Map(jen.String()).Qual(def.PkgCBOR, "RawMessage")
			g.If(
				jen.Err().Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "Unmarshal").Call(
					jen.Id("data"),
					jen.Op("&").Id("rawMap"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)

			// Unmarshal ID field for nodes and edges
			if isNode || isEdge {
				g.Line()
				g.Comment("Embedded som.Node/Edge ID field")
				g.If(
					jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit("id")),
					jen.Id("ok"),
				).BlockFunc(func(bg *jen.Group) {
					bg.Var().Id("recordID").Op("*").Qual(b.subPkg(""), "ID")
					bg.If(
					jen.Err().Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("recordID")),
					jen.Err().Op("!=").Nil(),
				).Block(jen.Return(jen.Err()))
					bg.Var().Id("idStr").String()
					bg.If(jen.Id("recordID").Op("!=").Nil()).Block(
						jen.List(jen.Id("s"), jen.Id("ok")).Op(":=").Id("recordID").Dot("ID").Assert(jen.String()),
						jen.If(jen.Op("!").Id("ok")).Block(
							jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("expected string ID, got %T"), jen.Id("recordID").Dot("ID"))),
						),
						jen.Id("idStr").Op("=").Id("s"),
					)

					if isNode {
						node := elem.(*field.NodeTable)
						fieldName := node.Source.EmbeddedFieldName
						if fieldName == "" {
							fieldName = "Node"
						}
						if fieldName == "Node" {
							bg.Id("c").Dot(fieldName).Op("=").Qual(b.subPkg(""), "NewNode").Call(jen.Id("idStr"))
						} else {
							bg.Id("c").Dot(fieldName).Op("=").Qual(b.subPkg(""), "NewCustomNode").Types(
								jen.Qual(b.subPkg(""), string(node.Source.IDType)),
							).Call(jen.Id("idStr"))
						}
					} else {
						bg.Id("c").Dot("Edge").Op("=").Qual(b.subPkg(""), "NewEdge").Call(jen.Id("idStr"))
					}
				})
			}

			// Unmarshal all fields
			g.Line()
			for _, f := range elem.GetFields() {
				// Skip ID field (handled specially for nodes/edges)
				if f.NameDatabase() == "id" {
					continue
				}

				// Generate unmarshal code for this field using field's CodeGen method
				if code := f.CodeGen().CBORUnmarshal(ctx); code != nil {
					g.Add(code)
				}
			}

			g.Line()
			g.Return(jen.Nil())
		})
}

func (b *convBuilder) buildTo(elem field.Element) jen.Code {
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
		// NO PTR - shallow wrapper: just unwrap
		jen.Func().
			Id(methodPrefix+elem.NameGo()).
			Params(jen.Id("data").Add(ptr).Id(localName)).
			Add(b.SourceQual(elem.NameGo())).
			Block(
				jen.Return(jen.Id("data").Dot(elem.NameGo())), // Just unwrap the embedding
			),

		jen.Line(),

		// PTR - shallow wrapper: just unwrap
		jen.Func().
			Id(methodPrefix+elem.NameGo()+"Ptr").
			Params(jen.Id("data").Op("*").Id(localName)).
			Op("*").Add(b.SourceQual(elem.NameGo())).
			Block(
				jen.If(jen.Id("data").Op("==").Nil()).Block(
					jen.Return(jen.Nil()),
				),

				jen.Id("result").Op(":=").Id("data").Dot(elem.NameGo()),
				jen.Return(jen.Op("&").Id("result")), // Unwrap and return pointer
			),
	)
}

func (b *convBuilder) buildFromLink(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"Link").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower()+"Link")).
		Add(b.SourceQual(node.NameGo())).
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("link").Op("==").Nil()).Block(
				jen.Return(jen.Add(b.SourceQual(node.NameGo())).Values()),
			)
			g.Id("res").Op(":=").Id(node.NameGo()).Call(jen.Id("link").Dot(node.NameGo()))
			g.Return(jen.Id("To" + node.NameGo()).Call(jen.Id("res")))
		})
}

func (b *convBuilder) buildFromLinkPtr(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"LinkPtr").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower()+"Link")).
		Op("*").Add(b.SourceQual(node.NameGo())).
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id("link").Op("==").Nil()).Block(
				jen.Return(jen.Nil()),
			)
			g.Id("res").Op(":=").Id(node.NameGo()).Call(jen.Id("link").Dot(node.NameGo()))
			g.Id("out").Op(":=").Id("To" + node.NameGo()).Call(jen.Id("res"))
			g.Return(jen.Op("&").Id("out"))
		})
}

func (b *convBuilder) buildToLink(node *field.NodeTable) jen.Code {
	tableName := node.NameDatabase()
	return jen.Func().
		Id("to"+node.NameGo()+"Link").
		Params(jen.Id("node").Add(b.SourceQual(node.NameGo()))).
		Op("*").Id(node.NameGoLower()+"Link").
		Block(
			jen.If(jen.Id("node").Dot("ID").Call().Op("==").Lit("")).Block(
				jen.Return(jen.Nil()),
			),
			jen.Id("rid").Op(":=").Qual("github.com/surrealdb/surrealdb.go/pkg/models", "NewRecordID").Call(
				jen.Lit(tableName), jen.Id("node").Dot("ID").Call(),
			),
			jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
				jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(jen.Id("node")),
				jen.Id("ID").Op(":").Op("&").Id("rid"),
			),
			jen.Return(jen.Op("&").Id("link")),
		)
}

func (b *convBuilder) buildToLinkPtr(node *field.NodeTable) jen.Code {
	tableName := node.NameDatabase()
	return jen.Func().
		Id("to"+node.NameGo()+"LinkPtr").
		Params(jen.Id("node").Op("*").Add(b.SourceQual(node.NameGo()))).
		Op("*").Id(node.NameGoLower()+"Link").
		Block(
			jen.
				If(
					jen.Id("node").Op("==").Nil().Op("||").
						Id("node").Dot("ID").Call().Op("==").Lit(""),
				).
				Block(
					jen.Return(jen.Nil()),
				),
			jen.Id("rid").Op(":=").Qual("github.com/surrealdb/surrealdb.go/pkg/models", "NewRecordID").Call(
				jen.Lit(tableName), jen.Id("node").Dot("ID").Call(),
			),
			jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
				jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(jen.Op("*").Id("node")),
				jen.Id("ID").Op(":").Op("&").Id("rid"),
			),
			jen.Return(jen.Op("&").Id("link")),
		)
}
