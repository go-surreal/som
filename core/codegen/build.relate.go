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
		{{range .EdgeFields}}
		func (n {{$.NameGo}}) {{.FieldName}}() {{.EdgeTypeLower}} {
			return {{.EdgeTypeLower}}(n)
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
			if edge.{{.InNameGo}}.ID() == "" {
				return errors.New("ID of the incoming node '{{.InNameGo}}' must not be empty")
			}
			if edge.{{.OutNameGo}}.ID() == "" {
				return errors.New("ID of the outgoing node '{{.OutNameGo}}' must not be empty")
			}
			inID := models.NewRecordID("{{.InNameDB}}", {{.InIDValue}})
			outID := models.NewRecordID("{{.OutNameDB}}", {{.OutIDValue}})
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
		"InNameGo":   edge.In.NameGo(),
		"InNameDB":   edge.In.Table().NameDatabase(),
		"InIDValue":  file.code(b.edgeNodeIDValue(edge.In)),
		"OutNameGo":  edge.Out.NameGo(),
		"OutNameDB":  edge.Out.Table().NameDatabase(),
		"OutIDValue": file.code(b.edgeNodeIDValue(edge.Out)),
	}

	return file.render(
		b.fs.Writer(path.Join(b.path(), edge.FileName())),
		"relateEdge", tmpl, data,
	)
}

// edgeNodeIDValue returns the expression for the record ID of the given node
// of an edge.
func (b *relateBuilder) edgeNodeIDValue(node *field.Node) jen.Code {
	id := jen.Id("edge").Dot(node.Table().NameGo()).Dot("ID").Call()

	if node.Table().Source.IDType == parser.IDTypeUUID {
		return jen.Qual(b.relativePkgPath(), "UUID").Call(id)
	}

	return id
}
