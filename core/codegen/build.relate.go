package codegen

import (
	"path"

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

type relateNodeEdgeField struct {
	FieldName     string
	EdgeTypeLower string
}

const relateNodeTmpl = `func New{{.NameGo}}(db Database) *{{.NameGo}} {
	return &{{.NameGo}}{db: db}
}

type {{.NameGo}} struct {
	db Database
}
{{range .EdgeFields}}
func (n {{$.NameGo}}) {{.FieldName}}() {{.EdgeTypeLower}} {
	return {{.EdgeTypeLower}}(n)
}
{{end}}`

func (b *relateBuilder) buildNodeFile(node *field.NodeTable) error {
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

	return renderGoFile(
		b.fs.Writer(path.Join(b.path(), node.FileName())),
		b.pkgName,
		"relateNode",
		relateNodeTmpl,
		data,
	)
}

const relateEdgeTmpl = `{{.ImportBlock}}

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
	err = e.db.Unmarshal(res, &rawResult)
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

func (b *relateBuilder) buildEdgeFile(edge *field.EdgeTable) error {
	needsSom := edge.In.Table().Source.IDType == parser.IDTypeUUID ||
		edge.Out.Table().Source.IDType == parser.IDTypeUUID

	imports := []goImport{
		{Path: "context"},
		{Path: "errors"},
		{Path: "fmt"},
		{Alias: "conv", Path: b.relativePkgPath(def.PkgConv)},
		{Alias: "internal", Path: b.relativePkgPath(def.PkgInternal)},
		{Alias: "model", Path: b.sourcePkgPath},
		{Alias: "models", Path: def.PkgModels},
	}

	if needsSom {
		imports = append(imports, goImport{Alias: "som", Path: b.relativePkgPath()})
	}

	data := map[string]any{
		"ImportBlock": formatImportBlock(imports),
		"TypeName":    edge.NameGoLower(),
		"EdgeNameGo":  edge.NameGo(),
		"EdgeNameDB":  edge.NameDatabase(),
		"InNameGo":    edge.In.NameGo(),
		"InNameDB":    edge.In.Table().NameDatabase(),
		"InIDValue":   edgeNodeIDValueStr(edge.In),
		"OutNameGo":   edge.Out.NameGo(),
		"OutNameDB":   edge.Out.Table().NameDatabase(),
		"OutIDValue":  edgeNodeIDValueStr(edge.Out),
	}

	return renderGoFile(
		b.fs.Writer(path.Join(b.path(), edge.FileName())),
		b.pkgName,
		"relateEdge",
		relateEdgeTmpl,
		data,
	)
}

func edgeNodeIDValueStr(node *field.Node) string {
	accessor := "edge." + node.Table().NameGo() + ".ID()"
	if node.Table().Source.IDType == parser.IDTypeUUID {
		return "som.UUID(" + accessor + ")"
	}
	return accessor
}
