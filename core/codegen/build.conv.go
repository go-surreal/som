package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
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
			jen.Id("ID").Op("*").Qual(def.PkgModels, "RecordID"),
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

func (b *convBuilder) nodeIDValue(node *field.NodeTable, varName string) jen.Code {
	if node.HasComplexID() {
		return b.complexNodeIDValue(node, varName)
	}
	if node.Source.IDType == parser.IDTypeUUID {
		return jen.Qual(b.subPkg(""), "UUID").Call(jen.Id(varName).Dot("ID").Call())
	}
	return jen.Id(varName).Dot("ID").Call()
}

func (b *convBuilder) complexNodeIDValue(node *field.NodeTable, varName string) jen.Code {
	cid := node.Source.ComplexID

	if cid.Kind == parser.IDTypeArray {
		var elems []jen.Code
		for _, sf := range cid.Fields {
			elems = append(elems, b.marshalFieldValue(sf, varName))
		}
		return jen.Index().Any().Values(elems...)
	}

	// Object ID: map[string]any{...}
	dict := jen.Dict{}
	for _, sf := range cid.Fields {
		dict[jen.Lit(sf.DBName)] = b.marshalFieldValue(sf, varName)
	}
	return jen.Map(jen.String()).Any().Values(dict)
}

func (b *convBuilder) unmarshalComplexID(g *jen.Group, node *field.NodeTable) {
	cid := node.Source.ComplexID
	cborPkg := path.Join(b.basePkg, "internal/cbor")

	g.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit("id")),
		jen.Id("ok"),
	).BlockFunc(func(bg *jen.Group) {
		bg.Var().Id("recordID").Op("*").Qual(def.PkgModels, "RecordID")
		bg.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("recordID")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))

		bg.If(jen.Id("recordID").Op("!=").Nil()).BlockFunc(func(inner *jen.Group) {
			// Re-marshal recordID.ID to raw CBOR bytes for typed unmarshal
			inner.List(jen.Id("idRaw"), jen.Err()).Op(":=").Qual(cborPkg, "Marshal").Call(jen.Id("recordID").Dot("ID"))
			inner.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
			if cid.Kind == parser.IDTypeArray {
				inner.Var().Id("rawArr").Index().Qual(def.PkgCBOR, "RawMessage")
				inner.If(
					jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawArr")),
					jen.Err().Op("!=").Nil(),
				).Block(jen.Return(jen.Err()))
				inner.If(jen.Len(jen.Id("rawArr")).Op(">=").Lit(len(cid.Fields))).BlockFunc(func(arrBlock *jen.Group) {
					arrBlock.Var().Id("key").Qual(b.sourcePkgPath, cid.StructName)

					for i, sf := range cid.Fields {
						arrBlock.Add(b.unmarshalFieldAssign("key", sf, jen.Id("rawArr").Index(jen.Lit(i)), cborPkg))
					}

					arrBlock.Id("c").Dot(node.Source.IDEmbed).Op("=").
						Qual(b.subPkg(""), "NewCustomNode").Types(
						jen.Qual(b.sourcePkgPath, cid.StructName),
					).Call(jen.Id("key"))
				})
			} else {
				inner.Var().Id("rawObj").Map(jen.String()).Qual(def.PkgCBOR, "RawMessage")
				inner.If(
					jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawObj")),
					jen.Err().Op("!=").Nil(),
				).Block(jen.Return(jen.Err()))
				inner.Var().Id("key").Qual(b.sourcePkgPath, cid.StructName)

				for _, sf := range cid.Fields {
					inner.Add(b.unmarshalFieldAssign("key", sf, jen.Id("rawObj").Index(jen.Lit(sf.DBName)), cborPkg))
				}

				inner.Id("c").Dot(node.Source.IDEmbed).Op("=").
					Qual(b.subPkg(""), "NewCustomNode").Types(
					jen.Qual(b.sourcePkgPath, cid.StructName),
				).Call(jen.Id("key"))
			}
		})
	})
}

func (b *convBuilder) unmarshalFieldAssign(keyVar string, sf parser.ComplexIDField, accessor jen.Code, cborPkg string) jen.Code {
	switch f := sf.Field.(type) {
	case *parser.FieldString, *parser.FieldNumeric, *parser.FieldBool:
		return jen.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(accessor, jen.Op("&").Id(keyVar).Dot(sf.Name)),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))

	case *parser.FieldTime:
		errVar := jen.Id(sf.Name + "Err")
		return jen.BlockFunc(func(g *jen.Group) {
			g.Var().Add(errVar).Error()
			g.List(jen.Id(keyVar).Dot(sf.Name), errVar).Op("=").
				Qual(cborPkg, "UnmarshalDateTime").Call(accessor)
			g.If(errVar.Clone().Op("!=").Nil()).Block(jen.Return(errVar.Clone()))
		})

	case *parser.FieldDuration:
		errVar := jen.Id(sf.Name + "Err")
		return jen.BlockFunc(func(g *jen.Group) {
			g.Var().Add(errVar).Error()
			g.List(jen.Id(keyVar).Dot(sf.Name), errVar).Op("=").
				Qual(cborPkg, "UnmarshalDuration").Call(accessor)
			g.If(errVar.Clone().Op("!=").Nil()).Block(jen.Return(errVar.Clone()))
		})

	case *parser.FieldUUID:
		var unmarshalFunc string
		switch f.Package {
		case parser.UUIDPackageGoogle:
			unmarshalFunc = "UnmarshalUUIDGoogle"
		case parser.UUIDPackageGofrs:
			unmarshalFunc = "UnmarshalUUIDGofrs"
		default:
			unmarshalFunc = "UnmarshalUUIDGoogle"
		}
		errVar := jen.Id(sf.Name + "Err")
		return jen.BlockFunc(func(g *jen.Group) {
			g.Var().Add(errVar).Error()
			g.List(jen.Id(keyVar).Dot(sf.Name), errVar).Op("=").
				Qual(cborPkg, unmarshalFunc).Call(accessor)
			g.If(errVar.Clone().Op("!=").Nil()).Block(jen.Return(errVar.Clone()))
		})

	case *parser.FieldNode:
		return b.unmarshalNodeRef(sf, f, accessor, cborPkg)

	default:
		return jen.Null()
	}
}

func (b *convBuilder) unmarshalNodeRef(sf parser.ComplexIDField, f *parser.FieldNode, accessor jen.Code, cborPkg string) jen.Code {
	refNode := b.findNodeByName(f.Node)
	if refNode == nil {
		return jen.Null()
	}

	return jen.BlockFunc(func(g *jen.Group) {
		g.Var().Id("rid").Op("*").Qual(def.PkgModels, "RecordID")
		g.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(accessor, jen.Op("&").Id("rid")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))
		g.If(jen.Id("rid").Op("!=").Nil()).BlockFunc(func(inner *jen.Group) {
			if !refNode.HasComplexID() {
				inner.List(jen.Id("idRaw"), jen.Err()).Op(":=").Qual(cborPkg, "Marshal").Call(jen.Id("rid").Dot("ID"))
				inner.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
				inner.Var().Id("idStr").String()
				inner.If(
					jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("idStr")),
					jen.Err().Op("!=").Nil(),
				).Block(jen.Return(jen.Err()))

				idEmbed := refNode.Source.IDEmbed
				if idEmbed == "" {
					idEmbed = "Node"
				}

				if idEmbed == "Node" {
					inner.Id("key").Dot(sf.Name).Op("=").Qual(b.sourcePkgPath, refNode.NameGo()).Values(jen.Dict{
						jen.Id(idEmbed): jen.Qual(b.subPkg(""), "NewNode").Call(jen.Id("idStr")),
					})
				} else {
					inner.Id("key").Dot(sf.Name).Op("=").Qual(b.sourcePkgPath, refNode.NameGo()).Values(jen.Dict{
						jen.Id(idEmbed): jen.Qual(b.subPkg(""), "NewCustomNode").Types(
							jen.Qual(b.subPkg(""), string(refNode.Source.IDType)),
						).Call(jen.Qual(b.subPkg(""), string(refNode.Source.IDType)).Call(jen.Id("idStr"))),
					})
				}
			} else {
				b.unmarshalNodeRefComplex(inner, sf, refNode, cborPkg)
			}
		})
	})
}

func (b *convBuilder) unmarshalNodeRefComplex(g *jen.Group, sf parser.ComplexIDField, refNode *field.NodeTable, cborPkg string) {
	cid := refNode.Source.ComplexID

	g.List(jen.Id("idRaw"), jen.Err()).Op(":=").Qual(cborPkg, "Marshal").Call(jen.Id("rid").Dot("ID"))
	g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))

	if cid.Kind == parser.IDTypeArray {
		g.Var().Id("rawArr").Index().Qual(def.PkgCBOR, "RawMessage")
		g.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawArr")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))
		g.If(jen.Len(jen.Id("rawArr")).Op(">=").Lit(len(cid.Fields))).BlockFunc(func(arrBlock *jen.Group) {
			arrBlock.Var().Id("innerKey").Qual(b.sourcePkgPath, cid.StructName)
			for i, innerSF := range cid.Fields {
				arrBlock.Add(b.unmarshalFieldAssign("innerKey", innerSF, jen.Id("rawArr").Index(jen.Lit(i)), cborPkg))
			}
			arrBlock.Id("key").Dot(sf.Name).Op("=").Qual(b.sourcePkgPath, refNode.NameGo()).Values(jen.Dict{
				jen.Id(refNode.Source.IDEmbed): jen.Qual(b.subPkg(""), "NewCustomNode").Types(
					jen.Qual(b.sourcePkgPath, cid.StructName),
				).Call(jen.Id("innerKey")),
			})
		})
	} else {
		g.Var().Id("rawObj").Map(jen.String()).Qual(def.PkgCBOR, "RawMessage")
		g.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawObj")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))
		g.BlockFunc(func(objBlock *jen.Group) {
			objBlock.Var().Id("innerKey").Qual(b.sourcePkgPath, cid.StructName)
			for _, innerSF := range cid.Fields {
				objBlock.Add(b.unmarshalFieldAssign("innerKey", innerSF, jen.Id("rawObj").Index(jen.Lit(innerSF.DBName)), cborPkg))
			}
			objBlock.Id("key").Dot(sf.Name).Op("=").Qual(b.sourcePkgPath, refNode.NameGo()).Values(jen.Dict{
				jen.Id(refNode.Source.IDEmbed): jen.Qual(b.subPkg(""), "NewCustomNode").Types(
					jen.Qual(b.sourcePkgPath, cid.StructName),
				).Call(jen.Id("innerKey")),
			})
		})
	}
}

func (b *convBuilder) marshalFieldValue(sf parser.ComplexIDField, varName string) jen.Code {
	accessor := jen.Id(varName).Dot("ID").Call().Dot(sf.Name)
	return fieldValueFrom(b.input, b.basePkg, sf, accessor)
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
				if node, ok := elem.(*field.NodeTable); !ok || !node.HasComplexID() {
					fieldCount++
				}
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

				if node, ok := elem.(*field.NodeTable); ok && node.HasComplexID() {
					// Complex IDs: no ID marshaling needed, sub-fields are populated from the record ID.
				} else {
					var idValue jen.Code
					if node, ok := elem.(*field.NodeTable); ok {
						idValue = b.nodeIDValue(node, "c")
					} else {
						idValue = jen.Id("c").Dot("ID").Call()
					}

					g.If(jen.Id("c").Dot("ID").Call().Op("!=").Lit("")).Block(
						jen.Id("data").Index(jen.Lit("id")).Op("=").Qual(def.PkgModels, "NewRecordID").Call(
							jen.Lit(tableName), idValue,
						),
					)
				}
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

				if node, ok := elem.(*field.NodeTable); ok && node.HasComplexID() {
					b.unmarshalComplexID(g, node)
				} else {
					g.If(
						jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit("id")),
						jen.Id("ok"),
					).BlockFunc(func(bg *jen.Group) {
						bg.Var().Id("recordID").Op("*").Qual(def.PkgModels, "RecordID")
						bg.If(
							jen.Err().Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("recordID")),
							jen.Err().Op("!=").Nil(),
						).Block(jen.Return(jen.Err()))
						bg.Var().Id("idStr").String()
						bg.If(jen.Id("recordID").Op("!=").Nil()).Block(
							jen.List(jen.Id("s"), jen.Err()).Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "RecordIDToString").Call(jen.Id("recordID").Dot("ID")),
							jen.If(jen.Err().Op("!=").Nil()).Block(
								jen.Return(jen.Err()),
							),
							jen.Id("idStr").Op("=").Id("s"),
						)

						if isNode {
							node := elem.(*field.NodeTable)
							fieldName := node.Source.IDEmbed
							if fieldName == "" {
								fieldName = "Node"
							}
							if fieldName == "Node" {
								bg.Id("c").Dot(fieldName).Op("=").Qual(b.subPkg(""), "NewNode").Call(jen.Id("idStr"))
							} else {
								bg.Id("c").Dot(fieldName).Op("=").Qual(b.subPkg(""), "NewCustomNode").Types(
									jen.Qual(b.subPkg(""), string(node.Source.IDType)),
								).Call(jen.Qual(b.subPkg(""), string(node.Source.IDType)).Call(jen.Id("idStr")))
							}
						} else {
							bg.Id("c").Dot("Edge").Op("=").Qual(b.subPkg(""), "NewEdge").Call(jen.Id("idStr"))
						}
					})
				}
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
		Id("from" + node.NameGo() + "Link").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower() + "Link")).
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
		Id("from" + node.NameGo() + "LinkPtr").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower() + "Link")).
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
	return b.buildToLinkCommon(node, false)
}

func (b *convBuilder) buildToLinkPtr(node *field.NodeTable) jen.Code {
	return b.buildToLinkCommon(node, true)
}

func (b *convBuilder) buildToLinkCommon(node *field.NodeTable, isPtr bool) jen.Code {
	tableName := node.NameDatabase()
	idVal := b.nodeIDValue(node, "node")

	var stmts []jen.Code

	if isPtr {
		stmts = append(stmts, jen.If(jen.Id("node").Op("==").Nil()).Block(
			jen.Return(jen.Nil()),
		))
	}

	if node.HasComplexID() {
		cid := node.Source.ComplexID
		if !cid.HasNodeRef() {
			stmts = append(stmts,
				jen.Var().Id("zeroKey").Add(b.SourceQual(cid.StructName)),
				jen.If(jen.Id("node").Dot("ID").Call().Op("==").Id("zeroKey")).Block(
					jen.Return(jen.Nil()),
				),
			)
		} else {
			b.addLinkNodeRefFieldChecks(&stmts, cid, "node")
		}
	} else {
		stmts = append(stmts, jen.If(jen.Id("node").Dot("ID").Call().Op("==").Lit("")).Block(
			jen.Return(jen.Nil()),
		))
	}

	var fromArg jen.Code
	if isPtr {
		fromArg = jen.Op("*").Id("node")
	} else {
		fromArg = jen.Id("node")
	}

	stmts = append(stmts,
		jen.Id("rid").Op(":=").Qual(def.PkgModels, "NewRecordID").Call(
			jen.Lit(tableName), idVal,
		),
		jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
			jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(fromArg),
			jen.Id("ID").Op(":").Op("&").Id("rid"),
		),
		jen.Return(jen.Op("&").Id("link")),
	)

	funcName := "to" + node.NameGo() + "Link"
	paramType := jen.Add(b.SourceQual(node.NameGo()))
	if isPtr {
		funcName += "Ptr"
		paramType = jen.Op("*").Add(b.SourceQual(node.NameGo()))
	}

	return jen.Func().
		Id(funcName).
		Params(jen.Id("node").Add(paramType)).
		Op("*").Id(node.NameGoLower() + "Link").
		Block(stmts...)
}

func (b *convBuilder) addLinkNodeRefFieldChecks(stmts *[]jen.Code, cid *parser.FieldComplexID, varName string) {
	for _, sf := range cid.Fields {
		fn, ok := sf.Field.(*parser.FieldNode)
		if !ok {
			continue
		}
		refNode := b.findNodeByName(fn.Node)
		if refNode == nil {
			continue
		}
		accessor := jen.Id(varName).Dot("ID").Call().Dot(sf.Name)
		if !refNode.HasComplexID() {
			*stmts = append(*stmts, jen.If(jen.Add(accessor).Dot("ID").Call().Op("==").Lit("")).Block(
				jen.Return(jen.Nil()),
			))
		} else if !refNode.Source.ComplexID.HasNodeRef() {
			zeroVar := "zero" + sf.Name + "Key"
			*stmts = append(*stmts,
				jen.Var().Id(zeroVar).Add(b.SourceQual(refNode.Source.ComplexID.StructName)),
				jen.If(jen.Add(accessor).Dot("ID").Call().Op("==").Id(zeroVar)).Block(
					jen.Return(jen.Nil()),
				),
			)
		}
	}
}
