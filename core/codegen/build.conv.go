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
				jen.Return(jen.Qual(def.PkgCBOR, "Marshal").Call(jen.Nil())),
			)

			// Count fields for pre-sized map allocation
			fieldCount := 0
			if isNode || isEdge {
				fieldCount++ // ID field
			}
			for _, f := range elem.GetFields() {
				if timeField, ok := f.(*field.Time); ok {
					if timeField.Source().IsCreatedAt || timeField.Source().IsUpdatedAt {
						fieldCount++ // Timestamp fields
						continue
					}
				}
				if f.NameDatabase() != "id" {
					fieldCount++ // Regular fields
				}
			}

			g.Id("data").Op(":=").Make(jen.Map(jen.String()).Any(), jen.Lit(fieldCount))

			// Marshal ID field for nodes and edges
			if isNode || isEdge {
				g.Line()
				g.Comment("Embedded som.Node/Edge ID field")
				g.If(jen.Id("c").Dot("ID").Call().Op("!=").Nil()).Block(
					jen.Id("data").Index(jen.Lit("id")).Op("=").Id("c").Dot("ID").Call(),
				)
			}

			// Marshal timestamp fields (CreatedAt, UpdatedAt)
			for _, f := range elem.GetFields() {
				if timeField, ok := f.(*field.Time); ok {
					if timeField.Source().IsCreatedAt || timeField.Source().IsUpdatedAt {
						g.Line()
						g.Comment("Embedded som.Timestamps field: " + f.NameGo())
						g.If(jen.Op("!").Id("c").Dot(f.NameGo()).Call().Dot("IsZero").Call()).Block(
							jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "DateTime").Values(
								jen.Id("Time").Op(":").Id("c").Dot(f.NameGo()).Call(),
							),
						)
					}
				}
			}

			// Marshal regular fields
			g.Line()
			g.Comment("Regular fields")
			for _, f := range elem.GetFields() {
				// Skip timestamp fields (already handled)
				if timeField, ok := f.(*field.Time); ok {
					if timeField.Source().IsCreatedAt || timeField.Source().IsUpdatedAt {
						continue
					}
				}

				// Skip ID field (handled specially for nodes/edges)
				if f.NameDatabase() == "id" {
					continue
				}

				// Generate marshal code for this field
				g.Add(b.buildFieldMarshal(f, ctx))
			}

			g.Line()
			g.Return(jen.Qual(def.PkgCBOR, "Marshal").Call(jen.Id("data")))
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
				jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(
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
					bg.Var().Id("id").Op("*").Qual(b.subPkg(""), "ID")
					bg.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("id"))

					if isNode {
						bg.Id("c").Dot("Node").Op("=").Qual(b.subPkg(""), "NewNode").Call(jen.Id("id"))
					} else {
						bg.Id("c").Dot("Edge").Op("=").Qual(b.subPkg(""), "NewEdge").Call(jen.Id("id"))
					}
				})
			}

			// Unmarshal timestamp fields
			var createdAtFound, updatedAtFound bool
			var createdAtVar, updatedAtVar string
			for _, f := range elem.GetFields() {
				if timeField, ok := f.(*field.Time); ok {
					if timeField.Source().IsCreatedAt {
						createdAtFound = true
						createdAtVar = "createdAt"
						g.Line()
						g.Comment("Embedded som.Timestamps field: CreatedAt")
						g.Var().Id(createdAtVar).Qual("time", "Time")
						g.If(
							jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit("created_at")),
							jen.Id("ok"),
						).BlockFunc(func(bg *jen.Group) {
							bg.Id(createdAtVar).Op(",").Id("_").Op("=").Qual(path.Join(b.basePkg, def.PkgCBORHelpers), "UnmarshalDateTime").Call(jen.Id("raw"))
						})
					} else if timeField.Source().IsUpdatedAt {
						updatedAtFound = true
						updatedAtVar = "updatedAt"
						g.Line()
						g.Comment("Embedded som.Timestamps field: UpdatedAt")
						g.Var().Id(updatedAtVar).Qual("time", "Time")
						g.If(
							jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit("updated_at")),
							jen.Id("ok"),
						).BlockFunc(func(bg *jen.Group) {
							bg.Id(updatedAtVar).Op(",").Id("_").Op("=").Qual(path.Join(b.basePkg, def.PkgCBORHelpers), "UnmarshalDateTime").Call(jen.Id("raw"))
						})
					}
				}
			}

			// Initialize Timestamps if we have timestamp fields
			if createdAtFound || updatedAtFound {
				g.Line()
				g.Comment("Initialize Timestamps embedding")
				createdPtr := jen.Nil()
				updatedPtr := jen.Nil()
				if createdAtFound {
					g.Id("createdAtDT").Op(":=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "DateTime").Values(
						jen.Dict{jen.Id("Time"): jen.Id(createdAtVar)},
					)
					createdPtr = jen.Id("createdAtDT")
				}
				if updatedAtFound {
					g.Id("updatedAtDT").Op(":=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "DateTime").Values(
						jen.Dict{jen.Id("Time"): jen.Id(updatedAtVar)},
					)
					updatedPtr = jen.Id("updatedAtDT")
				}
				g.Id("c").Dot("Timestamps").Op("=").Qual(b.subPkg(""), "NewTimestamps").Call(createdPtr, updatedPtr)
			}

			// Unmarshal regular fields
			g.Line()
			g.Comment("Regular fields")
			for _, f := range elem.GetFields() {
				// Skip timestamp fields
				if timeField, ok := f.(*field.Time); ok {
					if timeField.Source().IsCreatedAt || timeField.Source().IsUpdatedAt {
						continue
					}
				}

				// Skip ID field (handled specially for nodes/edges)
				if f.NameDatabase() == "id" {
					continue
				}

				// Generate unmarshal code for this field
				g.Add(b.buildFieldUnmarshal(f, ctx))
			}

			g.Line()
			g.Return(jen.Nil())
		})
}

func (b *convBuilder) buildFieldMarshal(f field.Field, ctx field.Context) jen.Code {
	// Helper function to generate marshal code for each field type
	// Check if field is a pointer in Go code
	isPointer := b.isFieldPointerInGo(f)

	// Wrap in nil check if pointer type
	if isPointer {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
			b.generateFieldMarshalCode(f, ctx, g)
		})
	}

	return jen.BlockFunc(func(g *jen.Group) {
		b.generateFieldMarshalCode(f, ctx, g)
	})
}

func (b *convBuilder) isFieldPointerInGo(f field.Field) bool {
	// Check specific field types that have Source() method with Pointer()
	switch tf := f.(type) {
	case *field.Time:
		return tf.Source().Pointer()
	case *field.Duration:
		return tf.Source().Pointer()
	case *field.UUID:
		return tf.Source().Pointer()
	case *field.URL:
		return tf.Source().Pointer()
	case *field.String:
		return tf.Source().Pointer()
	case *field.Bool:
		return tf.Source().Pointer()
	case *field.Numeric:
		return tf.Source().Pointer()
	case *field.Byte:
		return tf.Source().Pointer()
	case *field.Node:
		return tf.Source().Pointer()
	case *field.Struct:
		return tf.Source().Pointer()
	case *field.Slice:
		// Slices can be nil, so always check for nil before marshaling
		return true
	case *field.Edge, *field.Enum:
		// These types: For now assume they're not supported as pointers
		return false
	default:
		return false
	}
}

func (b *convBuilder) generateFieldMarshalCode(f field.Field, ctx field.Context, g *jen.Group) {
	switch tf := f.(type) {
	case *field.Time:
		// Direct assignment - types.DateTime has MarshalCBOR that cbor.Marshal will call
		if tf.Source().Pointer() {
			g.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
				jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "DateTime").Values(
					jen.Id("Time").Op(":").Op("*").Id("c").Dot(f.NameGo()),
				),
			)
		} else {
			g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "DateTime").Values(
				jen.Id("Time").Op(":").Id("c").Dot(f.NameGo()),
			)
		}
	case *field.Duration:
		// Direct assignment - types.Duration has MarshalCBOR that cbor.Marshal will call
		if tf.Source().Pointer() {
			g.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
				jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "Duration").Values(
					jen.Id("Duration").Op(":").Op("*").Id("c").Dot(f.NameGo()),
				),
			)
		} else {
			g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(b.basePkg, def.PkgTypes), "Duration").Values(
				jen.Id("Duration").Op(":").Id("c").Dot(f.NameGo()),
			)
		}
	case *field.UUID:
		// Direct assignment - types.UUID has MarshalCBOR that cbor.Marshal will call
		// UUID is a type alias, not a struct, so use type conversion
		if tf.Source().Pointer() {
			g.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).BlockFunc(func(bg *jen.Group) {
				bg.Id("uuidVal").Op(":=").Qual(path.Join(b.basePkg, def.PkgTypes), "UUID").Call(
					jen.Op("*").Id("c").Dot(f.NameGo()),
				)
				bg.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("uuidVal")
			})
		} else {
			g.Id("uuidVal").Op(":=").Qual(path.Join(b.basePkg, def.PkgTypes), "UUID").Call(
				jen.Id("c").Dot(f.NameGo()),
			)
			g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("uuidVal")
		}
	case *field.URL:
		// Convert URL to string for marshaling
		convFuncName := "fromURL"
		if tf.Source().Pointer() {
			convFuncName += "Ptr"
		}
		g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id(convFuncName).Call(jen.Id("c").Dot(f.NameGo()))
	case *field.Struct:
		// Convert to conv wrapper which has proper MarshalCBOR
		convFuncName := "from" + tf.Table().NameGo()
		if tf.Source().Pointer() {
			convFuncName += "Ptr"
		}
		g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id(convFuncName).Call(jen.Id("c").Dot(f.NameGo()))
	case *field.Node:
		// Node fields: convert to link (only ID, not full object)
		convFuncName := "to" + tf.Table().NameGo() + "Link"
		if tf.Source().Pointer() {
			convFuncName += "Ptr"
		}
		// For non-pointer fields, the conversion can still return nil if node has no ID
		// We need to check the result and only add if not nil
		if !tf.Source().Pointer() {
			g.If(jen.Id("link").Op(":=").Id(convFuncName).Call(jen.Id("c").Dot(f.NameGo())), jen.Id("link").Op("!=").Nil()).Block(
				jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("link"),
			)
		} else {
			g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id(convFuncName).Call(jen.Id("c").Dot(f.NameGo()))
		}
	default:
		// Simple types: direct assignment
		g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
	}
}

func (b *convBuilder) buildFieldUnmarshal(f field.Field, ctx field.Context) jen.Code {
	// Helper function to generate unmarshal code for each field type
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).BlockFunc(func(g *jen.Group) {
		switch tf := f.(type) {
		case *field.Time:
			// Use helper - check if pointer type
			helper := "UnmarshalDateTime"
			if tf.Source().Pointer() {
				helper = "UnmarshalDateTimePtr"
			}
			g.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(path.Join(b.basePkg, def.PkgCBORHelpers), helper).Call(jen.Id("raw"))
		case *field.Duration:
			// Use helper - check if pointer type
			helper := "UnmarshalDuration"
			if tf.Source().Pointer() {
				helper = "UnmarshalDurationPtr"
			}
			g.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(path.Join(b.basePkg, def.PkgCBORHelpers), helper).Call(jen.Id("raw"))
		case *field.UUID:
			// Use helper - check if pointer type
			helper := "UnmarshalUUID"
			if tf.Source().Pointer() {
				helper = "UnmarshalUUIDPtr"
			}
			g.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(path.Join(b.basePkg, def.PkgCBORHelpers), helper).Call(jen.Id("raw"))
		case *field.URL:
			// Convert string to URL for unmarshaling
			if tf.Source().Pointer() {
				g.Var().Id("convVal").Op("*").String()
				g.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
				g.Id("c").Dot(f.NameGo()).Op("=").Id("toURLPtr").Call(jen.Id("convVal"))
			} else {
				g.Var().Id("convVal").String()
				g.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
				g.Id("c").Dot(f.NameGo()).Op("=").Id("toURL").Call(jen.Id("convVal"))
			}
		case *field.Struct:
			// Unmarshal through conv wrapper
			if tf.Source().Pointer() {
				g.Var().Id("convVal").Op("*").Id(tf.Table().NameGoLower())
				g.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
				g.Id("c").Dot(f.NameGo()).Op("=").Id("to" + tf.Table().NameGo() + "Ptr").Call(jen.Id("convVal"))
			} else {
				g.Var().Id("convVal").Id(tf.Table().NameGoLower())
				g.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
				g.Id("c").Dot(f.NameGo()).Op("=").Id("to" + tf.Table().NameGo()).Call(jen.Id("convVal"))
			}
		case *field.Node:
			// Node fields: unmarshal through link (only ID, not full object)
			// convVal is always *groupLink, but we convert differently based on field type
			g.Var().Id("convVal").Op("*").Id(tf.Table().NameGoLower() + "Link")
			g.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
			if tf.Source().Pointer() {
				g.Id("c").Dot(f.NameGo()).Op("=").Id("from" + tf.Table().NameGo() + "LinkPtr").Call(jen.Id("convVal"))
			} else {
				g.Id("c").Dot(f.NameGo()).Op("=").Id("from" + tf.Table().NameGo() + "Link").Call(jen.Id("convVal"))
			}
		default:
			// Simple types: direct unmarshal
			g.Qual(def.PkgCBOR, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo()))
		}
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
