package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

type queryBuilder struct {
	*baseBuilder
}

func newQueryBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *queryBuilder {
	return &queryBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *queryBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	for _, view := range b.views {
		if err := b.buildViewFile(view); err != nil {
			return err
		}
	}

	return nil
}

func (b *queryBuilder) buildFile(node *field.NodeTable) error {
	hasRangeFn := node.HasComplexID() || node.HasStringID()

	file := b.newQueryFile()

	data := map[string]any{
		"NameGo":        node.NameGo(),
		"NameGoLower":   node.NameGoLower(),
		"NameDB":        node.NameDatabase(),
		"Kind":          "models",
		"HasFields":     true,
		"HasRangeFn":    hasRangeFn,
		"HasChangefeed": node.HasChangefeed(),
		"SoftDelete":    node.Source.SoftDelete,
		"Expiry":        node.Source.Expiry,

		"RangeBound": func(bound string) string { return file.code(b.rangeBound(node, bound)) },
	}

	return b.renderQueryFile(file, node.FileName(), data)
}

// buildViewFile generates the query builder constructor for a read-only view.
func (b *queryBuilder) buildViewFile(view *field.ViewTable) error {
	data := map[string]any{
		"NameGo":      view.NameGo(),
		"NameGoLower": view.NameGoLower(),
		"NameDB":      view.NameDatabase(),
		"Kind":        "views",
	}

	return b.renderQueryFile(b.newQueryFile(), view.FileName(), data)
}

func (b *queryBuilder) newQueryFile() *goFile {
	return newGoFile(b.pkgName,
		goImport{Alias: "som", Path: b.relativePkgPath()},
		goImport{Alias: "conv", Path: b.relativePkgPath(def.PkgConv)},
		goImport{Alias: "filter", Path: b.relativePkgPath(def.PkgFilter)},
		goImport{Alias: "lib", Path: b.relativePkgPath(def.PkgLib)},
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)
}

// renderQueryFile generates the model info and the query builder constructor
// for a single table. Views are read-only, so they have neither a soft-delete filter
// nor an ID range function.
//
// Deferred: SurrealDB's INCLUDE ORIGINAL (change-feed pre-image) is intentionally
// not implemented. It is a table-level clause (DEFINE TABLE ... CHANGEFEED ...
// INCLUDE ORIGINAL), not a SHOW CHANGES option, and enabling it changes the wire
// format: an update arrives as {current, update: [reverse json-patch incl. text
// diff]} instead of a full record. Proper support needs a tag opt-in, schema
// emission, and a diff-aware decoder that reverse-applies the patch to rebuild the
// original. Tracked separately; see the changefeed feature notes.
func (b *queryBuilder) renderQueryFile(file *goFile, fileName string, data map[string]any) error {
	tmpl := `
		// {{.NameGoLower}}ModelInfo holds the model-specific unmarshal functions for {{.NameGo}}.
		var {{.NameGoLower}}ModelInfo = modelInfo[model.{{.NameGo}}]{
			{{- if .HasFields}}
			Fields: conv.{{.NameGo}}Fields,
			{{- end}}
			UnmarshalAll: func(data []byte) ([]*model.{{.NameGo}}, error) {
				return unmarshalAll(data, conv.To{{.NameGo}}Ptr)
			},
			UnmarshalOne: func(data []byte) (*model.{{.NameGo}}, error) {
				return unmarshalOne(data, conv.To{{.NameGo}}Ptr)
			},
			UnmarshalSearchAll: func(data []byte, clauses []lib.SearchClause) ([]lib.SearchResult[*model.{{.NameGo}}], error) {
				return unmarshalSearchAll(data, clauses, conv.To{{.NameGo}}Ptr)
			},
		}
		{{if .HasRangeFn}}
		var {{.NameGoLower}}RangeFn = rangeFn[model.{{.NameGo}}](func(q *lib.Query[model.{{.NameGo}}], from som.RangeFrom, to som.RangeTo) string {
			expr := ":"
			if !from.IsOpen() {
				{{call .RangeBound "from"}}
			}
			if !from.IsOpen() && !from.IsInclusive() {
				expr += ">"
			}
			expr += ".."
			if !to.IsOpen() && to.IsInclusive() {
				expr += "="
			}
			if !to.IsOpen() {
				{{call .RangeBound "to"}}
			}
			return expr
		})
		{{end}}
		// New{{.NameGo}} creates a new query builder for {{.NameGo}} {{.Kind}}.
		func New{{.NameGo}}(db Database) Builder[model.{{.NameGo}}] {
			q := lib.NewQuery[model.{{.NameGo}}]("{{.NameDB}}")
			{{- if .SoftDelete}}
			// Automatically exclude soft-deleted records
			q.SoftDeleteFilter = filter.{{.NameGo}}.DeletedAt.Nil(true)
			{{- end}}
			{{- if .Expiry}}
			// Automatically exclude expired records
			q.ExpiryField = "expires_at"
			{{- end}}
			return Builder[model.{{.NameGo}}]{builder[model.{{.NameGo}}]{
				db:      db,
				info:    {{.NameGoLower}}ModelInfo,
				query:   q,
				{{- if .HasRangeFn}}
				rangeFn: {{.NameGoLower}}RangeFn,
				{{- end}}
			}}
		}
		{{if .HasChangefeed}}
		func New{{.NameGo}}Changes(db Database) ChangesBuilder[model.{{.NameGo}}, conv.{{.NameGo}}] {
			return ChangesBuilder[model.{{.NameGo}}, conv.{{.NameGo}}]{
				convFrom: conv.From{{.NameGo}}Ptr,
				convTo:   conv.To{{.NameGo}}Ptr,
				db:       db,
				table:    "{{.NameDB}}",
			}
		}
		{{end -}}
	`

	return file.render(
		b.fs.Writer(path.Join(b.path(), fileName)),
		"query", tmpl, data,
	)
}

// rangeBound emits the statements that append the given range bound to the
// record ID range expression.
func (b *queryBuilder) rangeBound(node *field.NodeTable, bound string) jen.Code {
	value := jen.Id(bound).Dot("Value").Call()
	expr := jen.Id("expr").Op("+=")

	if !node.HasComplexID() {
		idType := jen.Qual(b.relativePkgPath(), string(node.Source.IDType))

		return expr.Id("q").Dot("AsVar").Call(value.Assert(idType))
	}

	cid := node.Source.ComplexID
	keyType := b.SourceQual(cid.StructName)

	return jen.Id("key").Op(":=").Add(value.Assert(keyType)).Line().
		Add(expr.Add(b.rangeBoundExpr(node, cid, "key")))
}

func (b *queryBuilder) rangeBoundExpr(node *field.NodeTable, cid *parser.FieldComplexID, keyVar string) jen.Code {
	var parts []jen.Code

	if cid.Kind == parser.IDTypeArray {
		parts = append(parts, jen.Lit("["))
		for i, sf := range cid.Fields {
			if i > 0 {
				parts = append(parts, jen.Lit(", "))
			}
			parts = append(parts, b.rangeFieldAsVar(node, sf, keyVar))
		}
		parts = append(parts, jen.Lit("]"))
	} else {
		parts = append(parts, jen.Lit("{"))
		for i, sf := range cid.Fields {
			if i > 0 {
				parts = append(parts, jen.Lit(", "))
			}
			parts = append(parts, jen.Lit(sf.DBName+": "))
			parts = append(parts, b.rangeFieldAsVar(node, sf, keyVar))
		}
		parts = append(parts, jen.Lit("}"))
	}

	result := parts[0]
	for _, p := range parts[1:] {
		result = jen.Add(result).Op("+").Add(p)
	}
	return result
}

func (b *queryBuilder) rangeFieldAsVar(node *field.NodeTable, sf parser.ComplexIDField, keyVar string) jen.Code {
	accessor := jen.Id(keyVar).Dot(sf.Name)
	wrappedValue := fieldValueFrom(b.input, b.basePkg, sf, accessor)
	return jen.Id("q").Dot("AsVar").Call(wrappedValue)
}
