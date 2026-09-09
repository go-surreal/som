package codegen

import (
	"fmt"
	"path"
	"strings"

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

	for _, fragment := range b.fragments {
		if err := b.buildFragmentFile(fragment); err != nil {
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

// buildFragmentFile generates the CBOR conversion type of a fragment. A
// fragment is read-only and holds a subset of its parent's fields, so it only
// needs to decode - and to expose its projection for cursor values.
func (b *convBuilder) buildFragmentFile(fragment *field.FragmentTable) error {
	tmpl := `
		type {{.NameGo}} struct {
			model.{{.NameGo}}
		}

		func (c *{{.NameGo}}) fields() map[string]any {
			data := make(map[string]any, {{.FieldCount}})

			// Embedded som.Fragment record id
			if rid, ok := internal.FragmentRecordID(c.{{.NameGo}}); ok {
				data["id"] = rid
			}

			{{.MarshalFields}}

			return data
		}

		func (c *{{.NameGo}}) UnmarshalCBOR(data []byte) error {
			var rawMap map[string]cbor.RawMessage
			if err := cbor.Unmarshal(data, &rawMap); err != nil {
				return err
			}

			// Embedded som.Fragment record id
			if raw, ok := rawMap["id"]; ok {
				var recordID *models.RecordID
				if err := cbor.Unmarshal(raw, &recordID); err != nil {
					return err
				}
				internal.SetFragmentRecordID(&c.{{.NameGo}}.Fragment, recordID)
			}

			{{.UnmarshalFields}}

			// A fragment holds only a subset of the record's fields
			internal.SetMarker(&c.{{.NameGo}}.Fragment, internal.MarkerLoaded|internal.MarkerPartial)

			return nil
		}

		func To{{.NameGo}}(data {{.NameGo}}) model.{{.NameGo}} {
			return data.{{.NameGo}}
		}

		func To{{.NameGo}}Ptr(data *{{.NameGo}}) *model.{{.NameGo}} {
			if data == nil {
				return nil
			}
			result := data.{{.NameGo}}
			return &result
		}

		// {{.NameGo}}Fields returns the DB-keyed value map of the fragment, used to
		// derive pagination cursor values.
		func {{.NameGo}}Fields(m *model.{{.NameGo}}) map[string]any {
			c := {{.NameGo}}{*m}
			return c.fields()
		}
	`

	// The field code generation is bound to the table the fields belong to, so
	// that a fragment reuses the conversion helpers of its parent node.
	ctx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     fragment.Parent,
	}

	file := newGoFile(b.pkgName,
		goImport{Alias: "models", Path: def.PkgModels},
		goImport{Alias: "cbor", Path: b.relativePkgPath(def.PkgCBORHelpers)},
		goImport{Alias: "internal", Path: b.relativePkgPath(def.PkgInternal)},
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)

	data := map[string]any{
		"NameGo":     fragment.NameGo(),
		"FieldCount": len(fragment.Fields) + 1,

		"MarshalFields":   file.code(fieldCodes(fragment.Fields, ctx, (*field.CodeGen).CBORMarshal)),
		"UnmarshalFields": file.code(fieldCodes(fragment.Fields, ctx, (*field.CodeGen).CBORUnmarshal)),
	}

	return file.render(
		b.fs.Writer(path.Join(b.path(), fragment.FileName())),
		"convFragment", tmpl, data,
	)
}

// buildFile generates the CBOR conversion type for a single table or object.
// It is a shallow wrapper embedding the model, adding the CBOR marshalling
// that maps between Go field names and their database counterparts.
func (b *convBuilder) buildFile(elem field.Element) error {
	// Note: the space in "{ {{" is required, as "{{{" would start a template action.
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
			{{if .HasRelations}}
			// Flag the relations that hold no unresolved links
			{{.FetchedBitsCall}}
			{{end}}
			{{- if .UnmarshalsID}}
			// Mark the instance as fully loaded from the database
			internal.SetMarker(&c.{{.Embed}}, internal.MarkerLoaded)
			{{end}}
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
				// The link was not fetched, so only its record id is known.
				{{.LinkIDDecode}}
				internal.SetMarker(&f.{{.NameGo}}.{{.IDEmbed}}, internal.MarkerLoaded|internal.MarkerPartial)
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

		{{if .Relations}}
		// The relations of {{.NameGo}}, as bits of its load state.
		const (
			{{- range $i, $rel := .Relations}}
			{{$rel.Bit}}{{if not $i}} internal.Relations = 1 << iota{{end}}
			{{- end}}
		)

		// {{.NameGoLower}}FetchedBits reports which relations of the decoded record hold no
		// unresolved links. A relation without any link counts as resolved.
		func {{.NameGoLower}}FetchedBits(c *{{.NameGo}}) internal.Relations {
			var relations internal.Relations
			{{range $rel := .Relations}}
			{{$rel.Check}}
			{{end}}
			return relations
		}

		// {{.NameGo}}Resolved reports whether the given relation path of the model was
		// loaded from the database. The path is followed segment by segment, so a
		// nested path is only resolved if every relation along it was fetched.
		{{- if .Skipped}}
		//
		// The relation {{.Skipped}} is not tracked: its field nests links deeper than a
		// slice, so it is reported as not resolved and loaded again on every call.
		{{- end}}
		func {{.NameGo}}Resolved(m *model.{{.NameGo}}, path string) bool {
			head, rest, _ := strings.Cut(path, ".")
			switch head {
			{{- range $rel := .Relations}}
			case "{{$rel.NameDB}}":
				if !internal.Fetched(m).Has({{$rel.Bit}}) {
					return false
				}
				if rest == "" {
					return true
				}
				{{$rel.Nested}}
			{{- end}}
			}

			// An unknown relation is never resolved, so it is always fetched again.
			return false
		}
		{{else}}
		// {{.NameGo}}Resolved reports whether the given relation path of the model was
		// loaded from the database. The model has no relations, so no path is.
		{{- if .Skipped}}
		//
		// The relation {{.Skipped}} is not tracked: its field nests links deeper than a
		// slice, so it is reported as not resolved and loaded again on every call.
		{{- end}}
		func {{.NameGo}}Resolved(_ *model.{{.NameGo}}, _ string) bool {
			return false
		}
		{{end -}}
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

	embed := "View"
	switch {
	case isNode:
		embed = "Node"
	case isEdge:
		embed = "Edge"
	}

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
		goImport{Path: "strings"},
		goImport{Alias: "models", Path: def.PkgModels},
		goImport{Alias: "cbor", Path: b.relativePkgPath(def.PkgCBORHelpers)},
		goImport{Alias: "internal", Path: b.relativePkgPath(def.PkgInternal)},
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
		"Embed":       embed,
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
		if err := checkRelationLimit(node.NameGo(), len(relations(node))); err != nil {
			return err
		}

		data["Relations"] = b.relationData(file, node)
		data["Skipped"] = strings.Join(skippedRelations(node), ", ")
		data["HasRelations"] = len(relations(node)) > 0
		data["FetchedBitsCall"] = file.code(b.fetchedBitsCall(node))
		data["IDEmbed"] = node.Source.IDEmbed
		data["LinkIDChecks"] = file.code(b.linkIDChecks(node))
		data["LinkIDDecode"] = file.code(b.unmarshalLinkID(node))
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
	var fields []field.Field

	for _, f := range elem.GetFields() {
		if f.NameDatabase() == "id" {
			continue
		}

		fields = append(fields, f)
	}

	return fieldCodes(fields, ctx, codeOf)
}

// fieldCodes returns the given kind of conversion statement for all fields
// that need one.
func fieldCodes(
	fields []field.Field, ctx field.Context,
	codeOf func(*field.CodeGen, field.Context) jen.Code,
) jen.Code {
	var codes []jen.Code

	for _, f := range fields {
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

		bg.Add(b.unmarshalComplexIDInto(node, jen.Id("c")))
	})
}

// unmarshalLinkID returns the decoding of the record id of an unfetched link
// into the embedded som.Node of the link's model.
func (b *convBuilder) unmarshalLinkID(node *field.NodeTable) jen.Code {
	if node.HasComplexID() {
		return jen.Id("recordID").Op(":=").Id("f").Dot("ID").Line().
			Add(b.unmarshalComplexIDInto(node, jen.Id("f").Dot(node.NameGo())))
	}

	idType := jen.Qual(b.relativePkgPath(), string(node.Source.IDType))

	return jen.If(jen.Id("f").Dot("ID").Op("!=").Nil()).Block(
		jen.List(jen.Id("idStr"), jen.Err()).Op(":=").
			Qual(b.relativePkgPath(def.PkgCBORHelpers), "RecordIDToString").Call(jen.Id("f").Dot("ID").Dot("ID")),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		jen.Id("f").Dot(node.NameGo()).Dot(node.Source.IDEmbed).Op("=").
			Qual(b.relativePkgPath(), "NewNode").Types(idType).
			Call(idType.Clone().Call(jen.Id("idStr"))),
	)
}

// unmarshalComplexIDInto returns the decoding of a complex record id into the
// embedded som.Node of the given target. It expects a *models.RecordID variable
// named recordID to be in scope.
func (b *convBuilder) unmarshalComplexIDInto(node *field.NodeTable, target jen.Code) jen.Code {
	cid := node.Source.ComplexID
	cborPkg := path.Join(b.basePkg, "internal/cbor")

	newNode := jen.Qual(b.relativePkgPath(), "NewNode").
		Types(jen.Qual(b.sourcePkgPath, cid.StructName)).
		Call(jen.Id("key"))

	return jen.If(jen.Id("recordID").Op("!=").Nil()).BlockFunc(func(inner *jen.Group) {
		// Re-marshal recordID.ID to raw CBOR bytes for typed unmarshal
		inner.List(jen.Id("idRaw"), jen.Err()).Op(":=").Qual(cborPkg, "Marshal").Call(jen.Id("recordID").Dot("ID"))
		inner.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))

		if cid.Kind == parser.IDTypeArray {
			inner.Var().Id("rawArr").Index().Qual(cborPkg, "RawMessage")
			inner.If(
				jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawArr")),
				jen.Err().Op("!=").Nil(),
			).Block(jen.Return(jen.Err()))

			inner.If(jen.Len(jen.Id("rawArr")).Op(">=").Lit(len(cid.Fields))).BlockFunc(func(arrBlock *jen.Group) {
				arrBlock.Var().Id("key").Qual(b.sourcePkgPath, cid.StructName)

				for i, sf := range cid.Fields {
					arrBlock.Add(b.unmarshalFieldAssign("key", sf, jen.Id("rawArr").Index(jen.Lit(i)), cborPkg))
				}

				arrBlock.Add(target).Dot(node.Source.IDEmbed).Op("=").Add(newNode)
			})

			return
		}

		inner.Var().Id("rawObj").Map(jen.String()).Qual(cborPkg, "RawMessage")
		inner.If(
			jen.Err().Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("idRaw"), jen.Op("&").Id("rawObj")),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err()))
		inner.Var().Id("key").Qual(b.sourcePkgPath, cid.StructName)

		for _, sf := range cid.Fields {
			inner.Add(b.unmarshalFieldAssign("key", sf, jen.Id("rawObj").Index(jen.Lit(sf.DBName)), cborPkg))
		}

		inner.Add(target).Dot(node.Source.IDEmbed).Op("=").Add(newNode)
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

// relationData returns the per-relation values the conversion template needs:
// the bit of a relation, the check deriving it from a decoded record and the
// descent into the relation's target node.
func (b *convBuilder) relationData(file *goFile, node *field.NodeTable) []map[string]string {
	var out []map[string]string

	for _, rel := range relations(node) {
		out = append(out, map[string]string{
			"NameDB": rel.field.NameDatabase(),
			"Bit":    bitName(node, rel),
			"Check":  file.code(b.fetchedBit(node, rel)),
			"Nested": file.code(b.nestedResolved(rel)),
		})
	}

	return out
}

// fetchedBit returns the check that sets the bit of a single relation.
func (b *convBuilder) fetchedBit(node *field.NodeTable, rel relation) jen.Code {
	accessor := jen.Id("c").Dot(rel.field.NameGo())
	setBit := jen.Id("relations").Op("|=").Id(bitName(node, rel))

	if !rel.slice {
		if rel.elemPtr {
			return jen.If(jen.Add(accessor).Op("==").Nil().Op("||").Op("!").Add(accessor.Clone()).Dot("IsPartial").Call()).
				Block(setBit)
		}

		return jen.If(jen.Op("!").Add(accessor).Dot("IsPartial").Call()).Block(setBit)
	}

	// A slice is fetched as a whole, so it only counts as resolved if none of
	// its elements is still a partial link.
	return jen.BlockFunc(func(bg *jen.Group) {
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

// nestedResolved returns the descent into the target node of a relation, for
// the remaining segments of a path.
func (b *convBuilder) nestedResolved(rel relation) jen.Code {
	accessor := jen.Id("m").Dot(rel.field.NameGo())
	resolvedFn := rel.target.NameGo() + "Resolved"

	if !rel.slice {
		if rel.elemPtr {
			return joinStatements([]jen.Code{
				jen.If(jen.Add(accessor).Op("==").Nil()).Block(
					jen.Return(jen.True()),
				),
				jen.Return(jen.Id(resolvedFn).Call(accessor.Clone(), jen.Id("rest"))),
			})
		}

		return jen.Return(jen.Id(resolvedFn).Call(jen.Op("&").Add(accessor), jen.Id("rest")))
	}

	var stmts []jen.Code

	elems := accessor
	if rel.slicePtr {
		stmts = append(stmts, jen.If(jen.Add(accessor).Op("==").Nil()).Block(
			jen.Return(jen.True()),
		))
		elems = jen.Op("*").Add(accessor.Clone())
	}

	arg := jen.Id("v")
	if !rel.elemPtr {
		arg = jen.Op("&").Id("v")
	}

	stmts = append(stmts,
		jen.For(jen.List(jen.Id("_"), jen.Id("v")).Op(":=").Range().Add(elems)).BlockFunc(func(fg *jen.Group) {
			if rel.elemPtr {
				fg.If(jen.Id("v").Op("==").Nil()).Block(jen.Continue())
			}
			fg.If(jen.Op("!").Id(resolvedFn).Call(arg, jen.Id("rest"))).Block(
				jen.Return(jen.False()),
			)
		}),
		jen.Return(jen.True()),
	)

	return joinStatements(stmts)
}

// fetchedBitsCall returns the call flagging the resolved relations of a decoded
// record on its load state.
func (b *convBuilder) fetchedBitsCall(node *field.NodeTable) jen.Code {
	if len(relations(node)) < 1 {
		return jen.Null()
	}

	return jen.Qual(b.relativePkgPath("internal"), "AddFetched").Call(
		jen.Op("&").Id("c").Dot("Node"),
		jen.Id(node.NameGoLower()+"FetchedBits").Call(jen.Id("c")),
	)
}
