package codegen

import (
	"fmt"
	"path"
	"strings"

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

	for _, view := range b.views {
		if err := b.buildFile(view); err != nil {
			return err
		}
	}

	for _, sink := range b.sinks {
		if err := b.buildFile(sink); err != nil {
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
	_, isView := elem.(*field.ViewTable)
	_, isSink := elem.(*field.SinkTable)

	typeName := elem.NameGoLower()
	if isNode || isEdge || isView || isSink {
		typeName = elem.NameGo()
	}

	f.Line()
	f.Type().Id(typeName).Struct(
		jen.Add(b.SourceQual(elem.NameGo())),
	)

	f.Line()
	f.Add(b.buildMarshalCBOR(elem, typeName, fieldCtx, isNode, isEdge, isView))

	f.Line()
	f.Add(b.buildFields(elem, typeName, fieldCtx, isNode, isEdge))

	f.Line()
	f.Add(b.buildUnmarshalCBOR(elem, typeName, fieldCtx, isNode, isEdge, isView))

	f.Line()
	f.Add(b.buildFrom(elem))

	f.Line()
	f.Add(b.buildTo(elem))

	if node, ok := elem.(*field.NodeTable); ok {
		f.Line()
		f.Add(b.buildNodeFields(node, typeName))

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
				).BlockFunc(func(g *jen.Group) {
					g.Comment("The link was not fetched, so only its record id is known.")
					b.unmarshalLinkID(g, node)
					g.Add(b.setMarker(
						jen.Id("f").Dot(node.NameGo()),
						node.Source.IDEmbed,
						jen.Qual(b.relativePkgPath("internal"), "MarkerLoaded").
							Op("|").Qual(b.relativePkgPath("internal"), "MarkerPartial"),
					))
					g.Return(jen.Nil())
				}),

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

		if err := b.buildFetchedBits(f, node); err != nil {
			return err
		}
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
		return jen.Qual(b.relativePkgPath(), "UUID").Call(jen.Id(varName).Dot("ID").Call())
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

		b.unmarshalComplexIDInto(bg, node, jen.Id("c"))
	})
}

// unmarshalLinkID generates the decoding of the record id of an unfetched link
// into the embedded som.Node of the link's model.
func (b *convBuilder) unmarshalLinkID(g *jen.Group, node *field.NodeTable) {
	if node.HasComplexID() {
		g.Id("recordID").Op(":=").Id("f").Dot("ID")
		b.unmarshalComplexIDInto(g, node, jen.Id("f").Dot(node.NameGo()))
		return
	}

	idType := string(node.Source.IDType)

	g.If(jen.Id("f").Dot("ID").Op("!=").Nil()).Block(
		jen.List(jen.Id("idStr"), jen.Err()).Op(":=").
			Qual(path.Join(b.basePkg, "internal/cbor"), "RecordIDToString").Call(jen.Id("f").Dot("ID").Dot("ID")),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Id("f").Dot(node.NameGo()).Dot(node.Source.IDEmbed).Op("=").
			Qual(b.relativePkgPath(), "NewNode").Types(jen.Qual(b.relativePkgPath(), idType)).
			Call(jen.Qual(b.relativePkgPath(), idType).Call(jen.Id("idStr"))),
	)
}

// unmarshalComplexIDInto generates the decoding of a complex record id into the
// embedded som.Node of the given target. It expects a *models.RecordID variable
// named recordID to be in scope.
func (b *convBuilder) unmarshalComplexIDInto(g *jen.Group, node *field.NodeTable, target jen.Code) {
	cid := node.Source.ComplexID
	cborPkg := path.Join(b.basePkg, "internal/cbor")

	g.If(jen.Id("recordID").Op("!=").Nil()).BlockFunc(func(inner *jen.Group) {
		// Re-marshal recordID.ID to raw CBOR bytes for typed unmarshal
		inner.List(jen.Id("idRaw"), jen.Err()).Op(":=").Qual(cborPkg, "Marshal").Call(jen.Id("recordID").Dot("ID"))
		inner.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
		if cid.Kind == parser.IDTypeArray {
			inner.Var().Id("rawArr").Index().Qual(path.Join(b.basePkg, "internal/cbor"), "RawMessage")
			inner.If(
				jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawArr")),
				jen.Err().Op("!=").Nil(),
			).Block(jen.Return(jen.Err()))
			inner.If(jen.Len(jen.Id("rawArr")).Op(">=").Lit(len(cid.Fields))).BlockFunc(func(arrBlock *jen.Group) {
				arrBlock.Var().Id("key").Qual(b.sourcePkgPath, cid.StructName)

				for i, sf := range cid.Fields {
					arrBlock.Add(b.unmarshalFieldAssign("key", sf, jen.Id("rawArr").Index(jen.Lit(i)), cborPkg))
				}

				arrBlock.Add(target).Dot(node.Source.IDEmbed).Op("=").
					Qual(b.relativePkgPath(), "NewNode").Types(
					jen.Qual(b.sourcePkgPath, cid.StructName),
				).Call(jen.Id("key"))
			})
		} else {
			inner.Var().Id("rawObj").Map(jen.String()).Qual(path.Join(b.basePkg, "internal/cbor"), "RawMessage")
			inner.If(
				jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawObj")),
				jen.Err().Op("!=").Nil(),
			).Block(jen.Return(jen.Err()))
			inner.Var().Id("key").Qual(b.sourcePkgPath, cid.StructName)

			for _, sf := range cid.Fields {
				inner.Add(b.unmarshalFieldAssign("key", sf, jen.Id("rawObj").Index(jen.Lit(sf.DBName)), cborPkg))
			}

			inner.Add(target).Dot(node.Source.IDEmbed).Op("=").
				Qual(b.relativePkgPath(), "NewNode").Types(
				jen.Qual(b.sourcePkgPath, cid.StructName),
			).Call(jen.Id("key"))
		}
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
		case parser.UUIDPackageStd:
			unmarshalFunc = "UnmarshalUUIDStd"
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

				inner.Id("key").Dot(sf.Name).Op("=").Qual(b.sourcePkgPath, refNode.NameGo()).Values(jen.Dict{
					jen.Id("Node"): jen.Qual(b.relativePkgPath(), "NewNode").Types(
						jen.Qual(b.relativePkgPath(), string(refNode.Source.IDType)),
					).Call(jen.Qual(b.relativePkgPath(), string(refNode.Source.IDType)).Call(jen.Id("idStr"))),
				})
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
		g.Var().Id("rawArr").Index().Qual(path.Join(b.basePkg, "internal/cbor"), "RawMessage")
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
				jen.Id(refNode.Source.IDEmbed): jen.Qual(b.relativePkgPath(), "NewNode").Types(
					jen.Qual(b.sourcePkgPath, cid.StructName),
				).Call(jen.Id("innerKey")),
			})
		})
	} else {
		g.Var().Id("rawObj").Map(jen.String()).Qual(path.Join(b.basePkg, "internal/cbor"), "RawMessage")
		g.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawObj")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))
		g.Var().Id("innerKey").Qual(b.sourcePkgPath, cid.StructName)
		for _, innerSF := range cid.Fields {
			g.Add(b.unmarshalFieldAssign("innerKey", innerSF, jen.Id("rawObj").Index(jen.Lit(innerSF.DBName)), cborPkg))
		}
		g.Id("key").Dot(sf.Name).Op("=").Qual(b.sourcePkgPath, refNode.NameGo()).Values(jen.Dict{
			jen.Id(refNode.Source.IDEmbed): jen.Qual(b.relativePkgPath(), "NewNode").Types(
				jen.Qual(b.sourcePkgPath, cid.StructName),
			).Call(jen.Id("innerKey")),
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
	_, isView := elem.(*field.ViewTable)
	_, isSink := elem.(*field.SinkTable)

	if isNode || isEdge || isView || isSink {
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

func (b *convBuilder) buildMarshalCBOR(elem field.Element, typeName string, ctx field.Context, isNode, isEdge, isView bool) jen.Code {
	return jen.Func().
		Params(jen.Id("c").Op("*").Id(typeName)).
		Id("MarshalCBOR").Params().
		Params(jen.Index().Byte(), jen.Error()).
		Block(
			jen.If(jen.Id("c").Op("==").Nil()).Block(
				jen.Return(jen.Qual(path.Join(b.basePkg, "internal/cbor"), "Marshal").Call(jen.Nil())),
			),
			jen.Return(jen.Qual(path.Join(b.basePkg, "internal/cbor"), "Marshal").Call(jen.Id("c").Dot("fields").Call())),
		)
}

// buildFields generates a fields() method that builds the DB-keyed value map.
// It is used by MarshalCBOR and, for nodes, exposed via <Type>Fields for
// cursor-based pagination which needs DB field names and DB-typed values.
func (b *convBuilder) buildFields(elem field.Element, typeName string, ctx field.Context, isNode, isEdge bool) jen.Code {
	return jen.Func().
		Params(jen.Id("c").Op("*").Id(typeName)).
		Id("fields").Params().
		Map(jen.String()).Any().
		BlockFunc(func(g *jen.Group) {
			// Count fields for pre-sized map allocation.
			// Views are read-only and never marshal an id (their id may be a
			// composite that cannot be re-wrapped), so only nodes/edges count it.
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

			// Marshal ID field for nodes and edges. Views are read-only and
			// their id (possibly a composite GROUP BY key) is not marshaled.
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
			g.Return(jen.Id("data"))
		})
}

// buildNodeFields generates an exported <Type>Fields package function that
// returns the DB-keyed value map for a model. Used by the query builder to
// derive pagination cursor values with correct DB field names and types.
func (b *convBuilder) buildNodeFields(node *field.NodeTable, typeName string) jen.Code {
	return jen.Func().
		Id(node.NameGo() + "Fields").
		Params(jen.Id("m").Op("*").Add(b.SourceQual(node.NameGo()))).
		Map(jen.String()).Any().
		Block(
			jen.Id("c").Op(":=").Id(typeName).Values(jen.Op("*").Id("m")),
			jen.Return(jen.Id("c").Dot("fields").Call()),
		)
}

func (b *convBuilder) buildUnmarshalCBOR(elem field.Element, typeName string, ctx field.Context, isNode, isEdge, isView bool) jen.Code {
	return jen.Func().
		Params(jen.Id("c").Op("*").Id(typeName)).
		Id("UnmarshalCBOR").Params(jen.Id("data").Index().Byte()).
		Error().
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("rawMap").Map(jen.String()).Qual(path.Join(b.basePkg, "internal/cbor"), "RawMessage")
			g.If(
				jen.Err().Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "Unmarshal").Call(
					jen.Id("data"),
					jen.Op("&").Id("rawMap"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Err()),
			)

			// Unmarshal ID field for nodes, edges and views
			if isNode || isEdge || isView {
				g.Line()
				g.Comment("Embedded som.Node/Edge/View ID field")
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
						if isView {
							// A view's record id may be an array/object (e.g. a
							// GROUP BY composite key), so it is stored as the full
							// record-id string representation.
							bg.If(jen.Id("recordID").Op("!=").Nil()).Block(
								jen.Id("idStr").Op("=").Id("recordID").Dot("String").Call(),
							)
						} else {
							bg.If(jen.Id("recordID").Op("!=").Nil()).Block(
								jen.List(jen.Id("s"), jen.Err()).Op(":=").Qual(path.Join(b.basePkg, "internal/cbor"), "RecordIDToString").Call(jen.Id("recordID").Dot("ID")),
								jen.If(jen.Err().Op("!=").Nil()).Block(
									jen.Return(jen.Err()),
								),
								jen.Id("idStr").Op("=").Id("s"),
							)
						}

						if isNode {
							node := elem.(*field.NodeTable)
							bg.Id("c").Dot("Node").Op("=").Qual(b.relativePkgPath(), "NewNode").Types(
								jen.Qual(b.relativePkgPath(), string(node.Source.IDType)),
							).Call(jen.Qual(b.relativePkgPath(), string(node.Source.IDType)).Call(jen.Id("idStr")))
						} else if isView {
							bg.Id("c").Dot("View").Op("=").Qual(b.relativePkgPath(), "NewView").Call(jen.Id("idStr"))
						} else {
							bg.Id("c").Dot("Edge").Op("=").Qual(b.relativePkgPath(), "NewEdge").Call(jen.Id("idStr"))
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

			if node, ok := elem.(*field.NodeTable); ok {
				b.addFetchedBits(g, node)
			}

			if isNode || isEdge || isView {
				g.Line()
				g.Comment("Mark the instance as fully loaded from the database")
				g.Add(b.setMarker(jen.Id("c"), embedName(isNode, isEdge), jen.Qual(b.relativePkgPath("internal"), "MarkerLoaded")))
			}

			g.Line()
			g.Return(jen.Nil())
		})
}

func embedName(isNode, isEdge bool) string {
	switch {
	case isNode:
		return "Node"
	case isEdge:
		return "Edge"
	default:
		return "View"
	}
}

// setMarker generates the marker call for the embedded som.Node/Edge/View of
// the given receiver. The marker is set through the internal package, so that
// application code cannot change it.
func (b *convBuilder) setMarker(receiver jen.Code, embed string, flags jen.Code) jen.Code {
	return jen.Qual(b.relativePkgPath("internal"), "SetMarker").Call(
		jen.Op("&").Add(receiver).Dot(embed),
		flags,
	)
}

func (b *convBuilder) buildTo(elem field.Element) jen.Code {
	localName := elem.NameGoLower()
	methodPrefix := "to"

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)
	_, isView := elem.(*field.ViewTable)
	_, isSink := elem.(*field.SinkTable)

	if isNode || isEdge || isView || isSink {
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

// relation describes a relation field of a node, i.e. a field holding one or
// more links to another node. Relations are the fields that can be resolved
// through a FETCH clause.
type relation struct {
	field  field.Field
	target *field.NodeTable
	slice  bool

	// elemPtr and slicePtr hold whether the model field is a pointer to the
	// element and/or to the slice itself.
	elemPtr  bool
	slicePtr bool
}

// maxRelations is the number of relations a model may have. The load state
// holds them as bits of a uint64, hence the limit.
const maxRelations = 64

// checkRelationLimit reports whether the relations of a model still fit into
// its load state. Silently wrapping around would make Resolve skip a relation
// that was never loaded, so this is an error rather than a fallback.
func checkRelationLimit(name string, count int) error {
	if count > maxRelations {
		return fmt.Errorf(
			"model %s has %d relations, but at most %d are supported",
			name, count, maxRelations,
		)
	}

	return nil
}

// relations returns the relation fields of the given node, in the order that
// defines their bit position.
//
// Only a link and a flat slice of links are covered, each of them optionally
// behind a pointer. Any other shape holding links, e.g. a slice of slices, is
// left out on purpose: reaching its elements would need a nesting-aware walk
// for every level, which is not worth it for a shape that is exotic to begin
// with. A left-out field has no bit, so Resolved always reports false for it
// and Resolve fetches it again on every call. That is a lost optimisation, not
// a wrong result. To support such a shape, give it a bit here and handle its
// nesting in addFetchedBit and addNestedResolved.
func relations(node *field.NodeTable) []relation {
	var out []relation

	for _, fld := range node.GetFields() {
		switch typed := fld.(type) {

		case *field.Node:
			out = append(out, relation{
				field:   fld,
				target:  typed.Table(),
				elemPtr: typed.IsPointer(),
			})

		case *field.Slice:
			elem, ok := typed.Element().(*field.Node)
			if !ok {
				continue
			}
			out = append(out, relation{
				field:    fld,
				target:   elem.Table(),
				slice:    true,
				elemPtr:  elem.IsPointer(),
				slicePtr: typed.IsPointer(),
			})
		}
	}

	return out
}

// skippedRelations returns the fields of the given node that hold links, but
// whose shape is not covered by relations, so that the generated code can point
// them out.
func skippedRelations(node *field.NodeTable) []string {
	covered := make(map[string]bool)
	for _, rel := range relations(node) {
		covered[rel.field.NameGo()] = true
	}

	var out []string
	for _, fld := range node.GetFields() {
		if !covered[fld.NameGo()] && holdsNode(fld) {
			out = append(out, fld.NameGo())
		}
	}

	return out
}

// holdsNode reports whether the given field holds one or more links, at any
// level of nesting.
func holdsNode(fld field.Field) bool {
	switch typed := fld.(type) {
	case *field.Node:
		return true
	case *field.Slice:
		return holdsNode(typed.Element())
	}
	return false
}

// bitName returns the name of the constant holding the bit of a relation.
func bitName(node *field.NodeTable, rel relation) string {
	return node.NameGoLower() + "Fetched" + rel.field.NameGo()
}

// buildFetchedBits generates the bit constants of a node's relations, the
// function deriving them from a decoded record, and the function reporting
// whether a relation path is resolved.
func (b *convBuilder) buildFetchedBits(f *jen.File, node *field.NodeTable) error {
	rels := relations(node)

	if err := checkRelationLimit(node.NameGo(), len(rels)); err != nil {
		return err
	}


	skippedNote := func() {
		skipped := skippedRelations(node)
		if len(skipped) < 1 {
			return
		}
		f.Comment("")
		f.Commentf("The relation %s is not tracked: its field nests links deeper than a", strings.Join(skipped, ", "))
		f.Comment("slice, so it is reported as not resolved and loaded again on every call.")
	}

	// A node without relations still needs the Resolved function, as it may be
	// the target of a path that reaches further than the node can offer.
	if len(rels) < 1 {
		f.Line()
		f.Commentf("%sResolved reports whether the given relation path of the model was", node.NameGo())
		f.Comment("loaded from the database. The model has no relations, so no path is.")
		skippedNote()
		f.Func().Id(node.NameGo()+"Resolved").
			Params(
				jen.Id("_").Op("*").Add(b.SourceQual(node.NameGo())),
				jen.Id("_").String(),
			).
			Bool().
			Block(jen.Return(jen.False()))
		return nil
	}

	f.Line()
	f.Commentf("The relations of %s, as bits of its load state.", node.NameGo())
	f.Const().DefsFunc(func(g *jen.Group) {
		for i, rel := range rels {
			if i == 0 {
				g.Id(bitName(node, rel)).Uint64().Op("=").Lit(1).Op("<<").Iota()
				continue
			}
			g.Id(bitName(node, rel))
		}
	})

	pkgInternal := b.relativePkgPath("internal")

	f.Line()
	f.Commentf("%sFetchedBits reports which relations of the decoded record hold no", node.NameGoLower())
	f.Comment("unresolved links. A relation without any link counts as resolved.")
	f.Func().Id(node.NameGoLower()+"FetchedBits").
		Params(jen.Id("c").Op("*").Id(node.NameGo())).
		Uint64().
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("bits").Uint64()

			for _, rel := range rels {
				g.Line()
				b.addFetchedBit(g, node, rel)
			}

			g.Line()
			g.Return(jen.Id("bits"))
		})

	f.Line()
	f.Commentf("%sResolved reports whether the given relation path of the model was", node.NameGo())
	f.Comment("loaded from the database. The path is followed segment by segment, so a")
	f.Comment("nested path is only resolved if every relation along it was fetched.")
	skippedNote()
	f.Func().Id(node.NameGo()+"Resolved").
		Params(
			jen.Id("m").Op("*").Add(b.SourceQual(node.NameGo())),
			jen.Id("path").String(),
		).
		Bool().
		BlockFunc(func(g *jen.Group) {
			g.List(jen.Id("head"), jen.Id("rest"), jen.Id("_")).Op(":=").
				Qual("strings", "Cut").Call(jen.Id("path"), jen.Lit("."))

			g.Switch(jen.Id("head")).BlockFunc(func(sg *jen.Group) {
				for _, rel := range rels {
					sg.Case(jen.Lit(rel.field.NameDatabase())).BlockFunc(func(cg *jen.Group) {
						cg.If(
							jen.Qual(pkgInternal, "Fetched").Call(jen.Id("m")).
								Op("&").Id(bitName(node, rel)).Op("==").Lit(0),
						).Block(
							jen.Return(jen.False()),
						)
						cg.If(jen.Id("rest").Op("==").Lit("")).Block(
							jen.Return(jen.True()),
						)
						b.addNestedResolved(cg, rel)
					})
				}
			})

			g.Line()
			g.Comment("An unknown relation is never resolved, so it is always fetched again.")
			g.Return(jen.False())
		})

	return nil
}

// addFetchedBit generates the check that sets the bit of a single relation.
func (b *convBuilder) addFetchedBit(g *jen.Group, node *field.NodeTable, rel relation) {
	accessor := jen.Id("c").Dot(rel.field.NameGo())
	setBit := jen.Id("bits").Op("|=").Id(bitName(node, rel))

	if !rel.slice {
		if rel.elemPtr {
			g.If(jen.Add(accessor).Op("==").Nil().Op("||").Op("!").Add(accessor.Clone()).Dot("IsPartial").Call()).
				Block(setBit)
			return
		}

		g.If(jen.Op("!").Add(accessor).Dot("IsPartial").Call()).Block(setBit)
		return
	}

	// A slice is fetched as a whole, so it only counts as resolved if none of
	// its elements is still a partial link.
	g.BlockFunc(func(bg *jen.Group) {
		if rel.slicePtr {
			bg.If(jen.Add(accessor).Op("==").Nil()).Block(
				setBit,
			).Else().Block(
				jen.Id("resolved").Op(":=").True(),
				jen.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Op("*").Add(accessor.Clone())).Block(
					b.elemPartialCheck(rel),
				),
				jen.If(jen.Id("resolved")).Block(setBit.Clone()),
			)
			return
		}

		bg.Id("resolved").Op(":=").True()
		bg.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Add(accessor.Clone())).Block(
			b.elemPartialCheck(rel),
		)
		bg.If(jen.Id("resolved")).Block(setBit.Clone())
	})
}

// elemPartialCheck generates the loop body clearing the resolved flag for a
// partial slice element.
func (b *convBuilder) elemPartialCheck(rel relation) jen.Code {
	cond := jen.Id("v").Dot("IsPartial").Call()
	if rel.elemPtr {
		cond = jen.Id("v").Op("!=").Nil().Op("&&").Id("v").Dot("IsPartial").Call()
	}

	return jen.If(cond).Block(
		jen.Id("resolved").Op("=").False(),
		jen.Break(),
	)
}

// addNestedResolved generates the descent into the target node of a relation,
// for the remaining segments of a path.
func (b *convBuilder) addNestedResolved(g *jen.Group, rel relation) {
	accessor := jen.Id("m").Dot(rel.field.NameGo())
	resolvedFn := rel.target.NameGo() + "Resolved"

	if !rel.slice {
		if rel.elemPtr {
			g.If(jen.Add(accessor).Op("==").Nil()).Block(
				jen.Return(jen.True()),
			)
			g.Return(jen.Id(resolvedFn).Call(accessor.Clone(), jen.Id("rest")))
			return
		}

		g.Return(jen.Id(resolvedFn).Call(jen.Op("&").Add(accessor), jen.Id("rest")))
		return
	}

	elems := accessor
	if rel.slicePtr {
		g.If(jen.Add(accessor).Op("==").Nil()).Block(
			jen.Return(jen.True()),
		)
		elems = jen.Op("*").Add(accessor.Clone())
	}

	arg := jen.Id("v")
	if !rel.elemPtr {
		arg = jen.Op("&").Id("v")
	}

	g.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Add(elems)).BlockFunc(func(fg *jen.Group) {
		if rel.elemPtr {
			fg.If(jen.Id("v").Op("==").Nil()).Block(jen.Continue())
		}
		fg.If(jen.Op("!").Id(resolvedFn).Call(arg, jen.Id("rest"))).Block(
			jen.Return(jen.False()),
		)
	})
	g.Return(jen.True())
}

// addFetchedBits generates the call flagging the resolved relations of a
// decoded record on its load state.
func (b *convBuilder) addFetchedBits(g *jen.Group, node *field.NodeTable) {
	if len(relations(node)) < 1 {
		return
	}

	g.Line()
	g.Comment("Flag the relations that hold no unresolved links")
	g.Qual(b.relativePkgPath("internal"), "AddFetched").Call(
		jen.Op("&").Id("c").Dot("Node"),
		jen.Id(node.NameGoLower()+"FetchedBits").Call(jen.Id("c")),
	)
}
