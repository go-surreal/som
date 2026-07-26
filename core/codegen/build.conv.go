package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
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

// buildFile generates the CBOR conversion type for a single table or object.
// It is a shallow wrapper embedding the model, adding the CBOR marshalling
// that maps between Go field names and their database counterparts.
func (b *convBuilder) buildFile(elem field.Element) error {
	tmpl := `
		type {{.TypeName}} struct {
			model.{{.NameGo}}
		}

		func (c *{{.TypeName}}) MarshalCBOR() ([]byte, error) {
			if c == nil {
				return cbor.Marshal(nil)
			}
			return cbor.Marshal(c.fields())
		}

		func (c *{{.TypeName}}) fields() map[string]any {
			data := make(map[string]any, {{.FieldCount}})
			{{if .MarshalsID}}
			// Embedded som.Node/Edge ID field
			{{.MarshalID}}
			{{end}}
			{{.MarshalFields}}

			return data
		}

		func (c *{{.TypeName}}) UnmarshalCBOR(data []byte) error {
			var rawMap map[string]cbor.RawMessage
			if err := cbor.Unmarshal(data, &rawMap); err != nil {
				return err
			}
			{{if .UnmarshalsID}}
			// Embedded som.Node/Edge/View ID field
			{{.UnmarshalID}}
			{{end}}
			{{.UnmarshalFields}}

			return nil
		}

		func {{.FromPrefix}}{{.NameGo}}(data model.{{.NameGo}}) {{.TypeName}} {
			return {{.TypeName}}{ {{.NameGo}}: data}
		}
		func {{.FromPrefix}}{{.NameGo}}Ptr(data *model.{{.NameGo}}) *{{.TypeName}} {
			if data == nil {
				return nil
			}
			return &{{.TypeName}}{ {{.NameGo}}: *data}
		}

		func {{.ToPrefix}}{{.NameGo}}(data {{.ToParam}}) model.{{.NameGo}} {
			return data.{{.NameGo}}
		}
		func {{.ToPrefix}}{{.NameGo}}Ptr(data *{{.TypeName}}) *model.{{.NameGo}} {
			if data == nil {
				return nil
			}
			result := data.{{.NameGo}}
			return &result
		}
		{{if .IsNode}}
		// {{.NameGo}}Fields returns the database keyed value map of a model. It is used by
		// the query builder to derive pagination cursor values with correct database field
		// names and types.
		func {{.NameGo}}Fields(m *model.{{.NameGo}}) map[string]any {
			c := {{.TypeName}}{*m}
			return c.fields()
		}

		// {{.NameGoLower}}Link is a {{.NameGo}} as referenced by another record. It marshals
		// to its record ID only, but unmarshals from either a record ID or a fetched record.
		type {{.NameGoLower}}Link struct {
			{{.TypeName}}
			ID *models.RecordID
		}

		func (f *{{.NameGoLower}}Link) MarshalCBOR() ([]byte, error) {
			if f == nil {
				return nil, nil
			}
			return cbor.Marshal(f.ID)
		}

		func (f *{{.NameGoLower}}Link) UnmarshalCBOR(data []byte) error {
			if err := cbor.Unmarshal(data, &f.ID); err == nil {
				return nil
			}
			type alias {{.NameGoLower}}Link
			var link alias
			err := cbor.Unmarshal(data, &link)
			if err == nil {
				*f = {{.NameGoLower}}Link(link)
			}
			return err
		}

		func from{{.NameGo}}Link(link *{{.NameGoLower}}Link) model.{{.NameGo}} {
			if link == nil {
				return model.{{.NameGo}}{}
			}
			res := {{.TypeName}}(link.{{.NameGo}})
			return To{{.NameGo}}(res)
		}

		func from{{.NameGo}}LinkPtr(link *{{.NameGoLower}}Link) *model.{{.NameGo}} {
			if link == nil {
				return nil
			}
			res := {{.TypeName}}(link.{{.NameGo}})
			out := To{{.NameGo}}(res)
			return &out
		}

		func to{{.NameGo}}Link(node model.{{.NameGo}}) *{{.NameGoLower}}Link {
			{{.LinkIDChecks}}
			rid := models.NewRecordID("{{.NameDB}}", {{.LinkID}})
			link := {{.NameGoLower}}Link{ {{.TypeName}}: From{{.NameGo}}(node), ID: &rid}
			return &link
		}

		func to{{.NameGo}}LinkPtr(node *model.{{.NameGo}}) *{{.NameGoLower}}Link {
			if node == nil {
				return nil
			}
			{{.LinkIDChecks}}
			rid := models.NewRecordID("{{.NameDB}}", {{.LinkID}})
			link := {{.NameGoLower}}Link{ {{.TypeName}}: From{{.NameGo}}(*node), ID: &rid}
			return &link
		}
		{{end -}}
	`

	ctx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	node, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)
	_, isView := elem.(*field.ViewTable)
	_, isSink := elem.(*field.SinkTable)

	// Objects are not tables of their own, so their conversion type is unexported.
	isTable := isNode || isEdge || isView || isSink

	// The ID is marshalled for nodes and edges only. Views are read-only and their
	// ID may be a composite (e.g. a GROUP BY key) that cannot be re-wrapped, so it
	// is never written back. Sinks are write-only and never read back at all.
	marshalsID := isNode || isEdge
	unmarshalsID := isNode || isEdge || isView

	typeName := elem.NameGoLower()
	fromPrefix, toPrefix := "from", "to"

	if isTable {
		typeName = elem.NameGo()
		fromPrefix, toPrefix = "From", "To"
	}

	// Edges are only ever converted back by pointer.
	toParam := typeName
	if isEdge {
		toParam = "*" + typeName
	}

	file := newGoFile(b.pkgName,
		goImport{Alias: "models", Path: def.PkgModels},
		goImport{Alias: "cbor", Path: b.relativePkgPath(def.PkgCBORHelpers)},
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)

	data := map[string]any{
		"NameGo":      elem.NameGo(),
		"NameGoLower": elem.NameGoLower(),
		"NameDB":      elem.NameDatabase(),
		"TypeName":    typeName,
		"ToParam":     toParam,
		"FromPrefix":  fromPrefix,
		"ToPrefix":    toPrefix,
		"IsNode":      isNode,
		"FieldCount":  b.fieldCount(elem, isNode, isEdge),

		"MarshalsID":      marshalsID,
		"UnmarshalsID":    unmarshalsID,
		"MarshalFields":   file.code(b.fieldCodes(elem, ctx, (*field.CodeGen).CBORMarshal)),
		"UnmarshalFields": file.code(b.fieldCodes(elem, ctx, (*field.CodeGen).CBORUnmarshal)),
	}

	if marshalsID {
		data["MarshalID"] = file.code(b.marshalID(elem, isNode))
	}

	if unmarshalsID {
		data["UnmarshalID"] = file.code(b.unmarshalID(elem, isNode, isView))
	}

	if isNode {
		data["LinkIDChecks"] = file.code(b.linkIDChecks(node))
		data["LinkID"] = file.code(b.nodeIDValue(node, "node"))
	}

	return file.render(
		b.fs.Writer(path.Join(b.path(), elem.FileName())),
		"conv", tmpl, data,
	)
}

// fieldCount pre-sizes the value map of the fields method.
func (b *convBuilder) fieldCount(elem field.Element, isNode, isEdge bool) int {
	count := 0

	if isNode || isEdge {
		if node, ok := elem.(*field.NodeTable); !ok || !node.HasComplexID() {
			count++
		}
	}

	for _, f := range elem.GetFields() {
		if f.NameDatabase() != "id" {
			count++
		}
	}

	return count
}

// fieldCodes returns the marshal or unmarshal statements of all fields. The ID
// field is skipped, as it is handled separately for tables.
func (b *convBuilder) fieldCodes(
	elem field.Element, ctx field.Context,
	codeOf func(*field.CodeGen, field.Context) jen.Code,
) jen.Code {
	var codes []jen.Code

	for _, f := range elem.GetFields() {
		if f.NameDatabase() == "id" {
			continue
		}

		if code := codeOf(f.CodeGen(), ctx); code != nil {
			codes = append(codes, code)
		}
	}

	return joinStatements(codes)
}

// marshalID returns the statement writing the record ID of a node or edge.
// Complex IDs are not written back, as their sub-fields are populated from the
// record ID instead.
func (b *convBuilder) marshalID(elem field.Element, isNode bool) jen.Code {
	node, _ := elem.(*field.NodeTable)

	if isNode && node.HasComplexID() {
		return jen.Null()
	}

	var idValue jen.Code = jen.Id("c").Dot("ID").Call()
	if isNode {
		idValue = b.nodeIDValue(node, "c")
	}

	return jen.If(jen.Id("c").Dot("ID").Call().Op("!=").Lit("")).Block(
		jen.Id("data").Index(jen.Lit("id")).Op("=").Qual(def.PkgModels, "NewRecordID").Call(
			jen.Lit(elem.NameDatabase()), idValue,
		),
	)
}

// unmarshalID returns the statement reading the embedded ID of a table.
func (b *convBuilder) unmarshalID(elem field.Element, isNode, isView bool) jen.Code {
	cborPkg := b.relativePkgPath(def.PkgCBORHelpers)

	if node, ok := elem.(*field.NodeTable); ok && node.HasComplexID() {
		return b.unmarshalComplexID(node)
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit("id")),
		jen.Id("ok"),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("recordID").Op("*").Qual(def.PkgModels, "RecordID")
		g.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("recordID")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))
		g.Var().Id("idStr").String()

		if isView {
			// A view's record id may be an array/object (e.g. a GROUP BY composite
			// key), so it is stored as the full record-id string representation.
			g.If(jen.Id("recordID").Op("!=").Nil()).Block(
				jen.Id("idStr").Op("=").Id("recordID").Dot("String").Call(),
			)
		} else {
			g.If(jen.Id("recordID").Op("!=").Nil()).Block(
				jen.List(jen.Id("s"), jen.Err()).Op(":=").Qual(cborPkg, "RecordIDToString").Call(jen.Id("recordID").Dot("ID")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Err()),
				),
				jen.Id("idStr").Op("=").Id("s"),
			)
		}

		switch {
		case isNode:
			node := elem.(*field.NodeTable)
			idType := jen.Qual(b.relativePkgPath(), string(node.Source.IDType))

			g.Id("c").Dot("Node").Op("=").Qual(b.relativePkgPath(), "NewNode").
				Types(idType).Call(idType.Clone().Call(jen.Id("idStr")))

		case isView:
			g.Id("c").Dot("View").Op("=").Qual(b.relativePkgPath(), "NewView").Call(jen.Id("idStr"))

		default:
			g.Id("c").Dot("Edge").Op("=").Qual(b.relativePkgPath(), "NewEdge").Call(jen.Id("idStr"))
		}
	})
}

// linkIDChecks returns the guard clauses that keep a link nil as long as the
// ID of the referenced node is not set.
func (b *convBuilder) linkIDChecks(node *field.NodeTable) jen.Code {
	var stmts []jen.Code

	switch {
	case !node.HasComplexID():
		stmts = append(stmts, jen.If(jen.Id("node").Dot("ID").Call().Op("==").Lit("")).Block(
			jen.Return(jen.Nil()),
		))

	case !node.Source.ComplexID.HasNodeRef():
		stmts = append(stmts,
			jen.Var().Id("zeroKey").Add(b.SourceQual(node.Source.ComplexID.StructName)),
			jen.If(jen.Id("node").Dot("ID").Call().Op("==").Id("zeroKey")).Block(
				jen.Return(jen.Nil()),
			),
		)

	default:
		b.addLinkNodeRefFieldChecks(&stmts, node.Source.ComplexID, "node")
	}

	return joinStatements(stmts)
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

func (b *convBuilder) unmarshalComplexID(node *field.NodeTable) jen.Code {
	cid := node.Source.ComplexID
	cborPkg := path.Join(b.basePkg, "internal/cbor")

	return jen.If(
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

					arrBlock.Id("c").Dot(node.Source.IDEmbed).Op("=").
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

				inner.Id("c").Dot(node.Source.IDEmbed).Op("=").
					Qual(b.relativePkgPath(), "NewNode").Types(
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
