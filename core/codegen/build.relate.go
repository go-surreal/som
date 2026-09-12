package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

type relateBuilder struct {
	*baseBuilder
}

func newRelateBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *relateBuilder {
	return &relateBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *relateBuilder) build() error {
	for _, node := range b.nodes {
		if node.HasComplexID() {
			continue
		}
		if err := b.buildNodeFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildEdgeFile(edge); err != nil {
			return err
		}
	}

	return nil
}

// relateNodeEdgeField describes one of the edges a node can be related through.
type relateNodeEdgeField struct {
	FieldName     string
	EdgeTypeLower string
}

func (b *relateBuilder) buildNodeFile(node *field.NodeTable) error {
	tmpl := `
		func New{{.NameGo}}(db Database) *{{.NameGo}} {
			return &{{.NameGo}}{db: db}
		}

		type {{.NameGo}} struct {
			db Database
		}
		{{range $edge := .EdgeFields}}
		func (n {{$.NameGo}}) {{$edge.FieldName}}() {{$edge.EdgeTypeLower}} {
			return {{$edge.EdgeTypeLower}}(n)
		}
		{{end}}
	`

	var edgeFields []relateNodeEdgeField

	for _, fld := range node.GetFields() {
		slice, ok := fld.(*field.Slice)
		if !ok {
			continue
		}

		edgeElement, ok := slice.Element().(*field.Edge)
		if !ok {
			continue
		}

		edgeFields = append(edgeFields, relateNodeEdgeField{
			FieldName:     fld.NameGo(),
			EdgeTypeLower: edgeElement.Table().NameGoLower(),
		})
	}

	data := map[string]any{
		"NameGo":     node.NameGo(),
		"EdgeFields": edgeFields,
	}

	return newGoFile(b.pkgName).render(
		b.fs.Writer(path.Join(b.path(), node.FileName())),
		"relateNode", tmpl, data,
	)
}

func (b *relateBuilder) buildEdgeFile(edge *field.EdgeTable) error {
	tmpl := `
		type {{.TypeName}} struct {
			db Database
		}

		// Create creates a new edge between the given nodes.
		// Note: The ID type if both nodes must be a string or number for now.
		func (e {{.TypeName}}) Create(ctx context.Context, edge *model.{{.EdgeNameGo}}) error {
			if edge == nil {
				return errors.New("the given edge must not be nil")
			}
			if edge.ID() != "" {
				return errors.New("ID must not be set for an edge to be created")
			}
			{{.InIDStmts}}
			{{.OutIDStmts}}
			query := "RELATE $inID->{{.EdgeNameDB}}->$outID CONTENT $data"
			data := conv.From{{.EdgeNameGo}}(*edge)
			res, err := e.db.Query(ctx, query, map[string]any{"inID": inID, "outID": outID, "data": data})
			if err != nil {
				return fmt.Errorf("could not create relation: %w", err)
			}
			var rawResult []internal.QueryResult[conv.{{.EdgeNameGo}}]
			err = cbor.Unmarshal(res, &rawResult)
			if err != nil {
				return fmt.Errorf("could not unmarshal relation: %w", err)
			}
			if len(rawResult) < 1 || len(rawResult[0].Result) < 1 {
				return errors.New("no result returned for relation")
			}
			convEdge := &rawResult[0].Result[0]
			*edge = conv.To{{.EdgeNameGo}}(convEdge)
			return nil
		}

		func ({{.TypeName}}) Update(edge *model.{{.EdgeNameGo}}) error {
			// TODO: implement!
			return errors.New("not yet implemented")
		}

		func ({{.TypeName}}) Delete(edge *model.{{.EdgeNameGo}}) error {
			// TODO: implement!
			// https://surrealdb.com/docs/surrealdb/surrealql/statements/delete#deleting-graph-edges
			return errors.New("not yet implemented")
		}
	`

	file := newGoFile(b.pkgName,
		goImport{Path: "context"},
		goImport{Path: "errors"},
		goImport{Path: "fmt"},
		goImport{Alias: "models", Path: def.PkgModels},
		goImport{Alias: "cbor", Path: b.relativePkgPath(def.PkgCBORHelpers)},
		goImport{Alias: "conv", Path: b.relativePkgPath(def.PkgConv)},
		goImport{Alias: "internal", Path: b.relativePkgPath(def.PkgInternal)},
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)

	data := map[string]any{
		"TypeName":   edge.NameGoLower(),
		"EdgeNameGo": edge.NameGo(),
		"EdgeNameDB": edge.NameDatabase(),
		"InIDStmts":  file.code(b.edgeEndID(edge.In, "inID", "incoming")),
		"OutIDStmts": file.code(b.edgeEndID(edge.Out, "outID", "outgoing")),
	}

	return file.render(
		b.fs.Writer(path.Join(b.path(), edge.FileName())),
		"relateEdge", tmpl, data,
	)
}

// edgeEndID emits the statements building the record ID of one end of an edge
// into the given variable. A union end may point to any of its members, so its
// record ID is resolved from the concrete model it holds.
func (b *relateBuilder) edgeEndID(end field.EdgeEnd, varName, direction string) jen.Code {
	accessor := jen.Id("edge").Dot(end.NameGo())

	union, isUnion := end.(*field.Union)
	if isUnion {
		return jen.List(jen.Id(varName), jen.Err()).Op(":=").
			Qual(b.relativePkgPath(def.PkgConv), union.Union().NameGo()+"RecordID").Call(accessor).
			Line().
			If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Qual("fmt", "Errorf").Call(
				jen.Lit(direction+" node '"+end.NameGo()+"': %w"), jen.Err(),
			)),
		)
	}

	node := end.(*field.Node)

	id := jen.Add(accessor.Clone()).Dot("ID").Call()
	if node.Table().Source.IDType == parser.IDTypeUUID {
		id = jen.Qual(b.relativePkgPath(), "UUID").Call(id)
	}

	return jen.If(jen.Add(accessor.Clone()).Dot("ID").Call().Op("==").Lit("")).Block(
		jen.Return(jen.Qual("errors", "New").Call(
			jen.Lit("ID of the "+direction+" node '"+end.NameGo()+"' must not be empty"),
		)),
	).Line().
		Id(varName).Op(":=").Qual(def.PkgModels, "NewRecordID").
		Call(jen.Lit(node.Table().NameDatabase()), id)
}
