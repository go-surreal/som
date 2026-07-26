package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
)

const filenameSchema = "schema.surql"

// Statements built up from optional clauses are defined as templates, so that
// the shape of the resulting SurrealQL stays readable.
var (
	tableStmt = template.Must(template.New("table").Parse(
		`DEFINE TABLE {{.Name}}{{if .Drop}} DROP{{end}} SCHEMAFULL TYPE {{.Type}}` +
			`{{if .In}} IN {{.In}} OUT {{.Out}} ENFORCED{{end}}` +
			`{{if .Changefeed}} CHANGEFEED {{.Changefeed}}{{end}} PERMISSIONS FULL;`,
	))

	viewStmt = template.Must(template.New("view").Parse(
		`DEFINE TABLE {{.Name}} TYPE NORMAL AS SELECT {{.Projections}} FROM {{.Source}}` +
			`{{if .Where}} WHERE {{.Where}}{{end}}` +
			`{{if .GroupBy}} GROUP BY {{.GroupBy}}{{end}};`,
	))

	searchIndexStmt = template.Must(template.New("searchIndex").Parse(
		`DEFINE INDEX {{.Name}} ON {{.Table}} FIELDS {{.Field}} FULLTEXT ANALYZER {{.Analyzer}}` +
			`{{if .BM25}} BM25({{.BM25K1}}, {{.BM25B}}){{else}} BM25{{end}}` +
			`{{if .Highlights}} HIGHLIGHTS{{end}}` +
			`{{if .Concurrently}} CONCURRENTLY{{end}};`,
	))

	analyzerStmt = template.Must(template.New("analyzer").Parse(
		`DEFINE ANALYZER {{.Name}}` +
			`{{if .Tokenizers}} TOKENIZERS {{.Tokenizers}}{{end}}` +
			`{{if .Filters}} FILTERS {{.Filters}}{{end}};`,
	))
)

// tableDef holds the parts of a DEFINE TABLE statement. In and Out are only
// set for relation tables, Drop only for write-only sinks.
type tableDef struct {
	Name       string
	Type       string
	Drop       bool
	In         string
	Out        string
	Changefeed string
}

func (b *build) buildSchemaFile() error {
	statements := []string{string(embed.CodegenComment), ""}

	// Generate DEFINE ANALYZER statements first
	if b.input.define != nil {
		for _, analyzer := range b.input.define.Analyzers {
			statement, err := buildAnalyzerStatement(analyzer)
			if err != nil {
				return err
			}

			statements = append(statements, statement)
		}
		if len(b.input.define.Analyzers) > 0 {
			statements = append(statements, "")
		}
	}

	// Collect index statements to add after table/field definitions
	var indexStatements []string

	for _, node := range b.input.nodes {
		statement, err := renderStatement(tableStmt, tableDef{
			Name:       node.NameDatabase(),
			Type:       "NORMAL",
			Changefeed: node.Changefeed,
		})
		if err != nil {
			return err
		}

		statements = append(statements, statement)

		for _, f := range node.GetFields() {
			statements = append(statements, f.SchemaStatements(node.NameDatabase(), "")...)
		}

		// Build indexes for this table (handles both simple and composite)
		tableIndexes, err := b.buildTableIndexStatements(node.NameDatabase(), node.GetFields(), node.Source.SoftDelete)
		if err != nil {
			return err
		}

		indexStatements = append(indexStatements, tableIndexes...)

		// Index expires_at to keep expiry purge deletes efficient.
		if node.Source.Expiry {
			indexName := fmt.Sprintf(def.IndexPrefix+"%s_expires_at", node.NameDatabase())
			indexStatements = append(indexStatements,
				fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS expires_at CONCURRENTLY;", indexName, node.NameDatabase()))
		}

		statements = append(statements, "")
	}

	for _, edge := range b.input.edges {
		statement, err := renderStatement(tableStmt, tableDef{
			Name:       edge.NameDatabase(),
			Type:       "RELATION",
			In:         edge.In.NameDatabase(),
			Out:        edge.Out.NameDatabase(), // TODO: can be OR'ed with "|"
			Changefeed: edge.Changefeed,
		})
		if err != nil {
			return err
		}

		statements = append(statements, statement)

		for _, f := range edge.GetFields() {
			statements = append(statements, f.SchemaStatements(edge.NameDatabase(), "")...)
		}

		// Build indexes for this table (handles both simple and composite)
		tableIndexes, err := b.buildTableIndexStatements(edge.NameDatabase(), edge.GetFields(), edge.Source.SoftDelete)
		if err != nil {
			return err
		}

		indexStatements = append(indexStatements, tableIndexes...)

		statements = append(statements, "")
	}

	// Sinks are write-only ingestion tables: records are accepted (firing
	// any dependent views/events) but discarded via DROP. Fields are still
	// defined so writes are validated and dependent view SELECTs typecheck.
	// No indexes are emitted — there are no rows to index.
	for _, sink := range b.input.sinks {
		statement, err := renderStatement(tableStmt, tableDef{
			Name: sink.NameDatabase(),
			Type: "NORMAL",
			Drop: true,
		})
		if err != nil {
			return err
		}

		statements = append(statements, statement)

		for _, f := range sink.GetFields() {
			statements = append(statements, f.SchemaStatements(sink.NameDatabase(), "")...)
		}

		statements = append(statements, "")
	}

	// Views are read-only, pre-computed tables defined via AS SELECT.
	// They are emitted after the source tables they depend on.
	for _, view := range b.input.views {
		statement, err := b.buildViewStatement(view)
		if err != nil {
			return err
		}
		if statement == "" {
			continue // view without a definition (see buildViewStatement)
		}
		statements = append(statements, statement, "")
	}

	// Append index statements at the end
	if len(indexStatements) > 0 {
		statements = append(statements, indexStatements...)
		statements = append(statements, "")
	}

	content := strings.Join(statements, "\n")

	b.fs.Write(path.Join(def.PkgRepo, "schema", filenameSchema), []byte(content))

	return nil
}

// buildViewStatement builds the DEFINE TABLE ... AS SELECT statement for a
// read-only view, joining the view model (its projected columns) with the
// SELECT definition supplied via a //go:build som definition file.
func (b *build) buildViewStatement(view *field.ViewTable) (string, error) {
	var def *parser.ViewDef
	if b.input.define != nil {
		for i := range b.input.define.Views {
			v := &b.input.define.Views[i]
			if v.View != view.NameGo() {
				continue
			}
			if def != nil {
				return "", fmt.Errorf(
					"view %s: multiple definitions found; multi-source views are not yet supported (SurrealDB #5593)",
					view.NameGo(),
				)
			}
			def = v
		}
	}

	if def == nil {
		// A view struct with no definition yet is not fatal: it lets the
		// read stack be generated first, so define.View can reference the
		// view's own filter refs, then a second gen emits the DDL. Warn and
		// skip emitting a statement for this view.
		fmt.Fprintf(os.Stderr,
			"warning: view %s has no definition; skipping its schema statement. "+
				"Declare it via define.View in a //go:build som file, then regenerate.\n",
			view.NameGo(),
		)
		return "", nil
	}

	if len(def.Projections) == 0 {
		return "", fmt.Errorf("view %s: definition has no projections", view.NameGo())
	}

	// A view may select from a node, an edge (relation) or a write-only
	// sink table (the sink→view ingestion pattern).
	var sourceDB string
	if node := b.input.findNodeByName(def.Source); node != nil {
		sourceDB = node.NameDatabase()
	} else if edge := b.input.findEdgeByName(def.Source); edge != nil {
		sourceDB = edge.NameDatabase()
	} else if sink := b.input.findSinkByName(def.Source); sink != nil {
		sourceDB = sink.NameDatabase()
	} else {
		return "", fmt.Errorf("view %s: unknown source model %q", view.NameGo(), def.Source)
	}

	return renderStatement(viewStmt, map[string]any{
		"Name":        view.NameDatabase(),
		"Projections": strings.Join(def.Projections, ", "),
		"Source":      sourceDB,
		"Where":       def.Where,
		"GroupBy":     strings.Join(def.GroupBy, ", "),
	})
}

func buildAnalyzerStatement(analyzer parser.AnalyzerDef) (string, error) {
	var filters []string
	for _, filter := range analyzer.Filters {
		filters = append(filters, buildFilterString(filter))
	}

	return renderStatement(analyzerStmt, map[string]any{
		"Name":       analyzer.Name,
		"Tokenizers": strings.Join(analyzer.Tokenizers, ", "),
		"Filters":    strings.Join(filters, ", "),
	})
}

func buildFilterString(filter parser.FilterDef) string {
	if len(filter.Params) == 0 {
		return filter.Name
	}

	var paramStrs []string
	for _, p := range filter.Params {
		switch v := p.(type) {
		case string:
			// Language identifiers (e.g., snowball) must NOT be quoted
			paramStrs = append(paramStrs, v)
		case int:
			paramStrs = append(paramStrs, fmt.Sprintf("%d", v))
		case float64:
			paramStrs = append(paramStrs, fmt.Sprintf("%g", v))
		default:
			paramStrs = append(paramStrs, fmt.Sprintf("%v", v))
		}
	}
	return fmt.Sprintf("%s(%s)", filter.Name, strings.Join(paramStrs, ", "))
}

// buildTableIndexStatements builds all index statements for a table, handling both
// simple indexes and composite unique indexes (fields grouped by name).
func (b *build) buildTableIndexStatements(tableName string, fields []field.Field, softDelete bool) ([]string, error) {
	var statements []string

	if !b.noCountIndex {
		stmt := fmt.Sprintf("DEFINE INDEX "+def.IndexPrefix+"%s_count ON %s COUNT;", tableName, tableName)
		statements = append(statements, stmt)
	}

	// Collect composite unique index fields grouped by name
	compositeUnique := make(map[string][]string) // name -> []fieldPath

	// Process all fields (including nested)
	if err := b.collectIndexes(tableName, "", fields, &statements, compositeUnique); err != nil {
		return nil, err
	}

	if softDelete {
		indexName := fmt.Sprintf(def.IndexPrefix+"%s_deleted_at", tableName)
		stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS deleted_at CONCURRENTLY;", indexName, tableName)
		statements = append(statements, stmt)
	}

	// Generate composite unique index statements
	for uniqueName, fieldPaths := range compositeUnique {
		// Index name format: __som__<table>_unique_<name>
		indexName := fmt.Sprintf(def.IndexPrefix+"%s_unique_%s", tableName, uniqueName)
		fieldsStr := strings.Join(fieldPaths, ", ")
		stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s UNIQUE;", indexName, tableName, fieldsStr)
		statements = append(statements, stmt)
	}

	return statements, nil
}

// collectIndexes recursively collects index statements and composite unique fields.
func (b *build) collectIndexes(tableName, fieldPrefix string, fields []field.Field, statements *[]string, compositeUnique map[string][]string) error {
	for _, f := range fields {
		fieldPath := f.NameDatabase()
		if fieldPrefix != "" {
			fieldPath = fieldPrefix + "." + fieldPath
		}

		for _, indexInfo := range f.Indexes() {
			if indexInfo.Unique && indexInfo.Name != "" {
				// Composite unique index - collect field for later
				compositeUnique[indexInfo.Name] = append(compositeUnique[indexInfo.Name], fieldPath)
			} else if indexInfo.Unique {
				// Simple unique index on single field
				indexName := fmt.Sprintf(def.IndexPrefix+"%s_unique_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s UNIQUE;", indexName, tableName, fieldPath)
				*statements = append(*statements, stmt)
			} else {
				// Regular (non-unique) index
				indexName := indexInfo.Name
				if indexName == "" {
					indexName = fmt.Sprintf(def.IndexPrefix+"%s_index_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				}
				stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s CONCURRENTLY;", indexName, tableName, fieldPath)
				*statements = append(*statements, stmt)
			}
		}

		searchInfo := f.SearchInfo()
		if searchInfo != nil && searchInfo.ConfigName != "" {
			// Look up the search config to get analyzer and options
			searchDef := b.findSearchConfig(searchInfo.ConfigName)
			if searchDef != nil {
				// Index name format: __som__<table>_search_<field>
				indexName := fmt.Sprintf(def.IndexPrefix+"%s_search_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))

				stmt, err := renderStatement(searchIndexStmt, map[string]any{
					"Name":         indexName,
					"Table":        tableName,
					"Field":        fieldPath,
					"Analyzer":     searchDef.AnalyzerName,
					"BM25":         searchDef.HasBM25,
					"BM25K1":       searchDef.BM25K1,
					"BM25B":        searchDef.BM25B,
					"Highlights":   searchDef.Highlights,
					"Concurrently": searchDef.Concurrently,
				})
				if err != nil {
					return err
				}

				*statements = append(*statements, stmt)
			}
		}

		// Handle nested struct fields
		if nestedFields := f.NestedFields(); nestedFields != nil {
			if err := b.collectIndexes(tableName, fieldPath, nestedFields, statements, compositeUnique); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *build) findSearchConfig(name string) *parser.SearchDef {
	if b.input.define == nil {
		return nil
	}
	for i := range b.input.define.Searches {
		if b.input.define.Searches[i].Name == name {
			return &b.input.define.Searches[i]
		}
	}
	return nil
}
