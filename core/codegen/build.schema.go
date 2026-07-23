package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
)

const filenameSchema = "schema.surql"

func (b *build) buildSchemaFile() error {
	statements := []string{string(embed.CodegenComment), ""}

	// Generate DEFINE ANALYZER statements first
	if b.input.define != nil {
		for _, analyzer := range b.input.define.Analyzers {
			statements = append(statements, buildAnalyzerStatement(analyzer))
		}
		if len(b.input.define.Analyzers) > 0 {
			statements = append(statements, "")
		}
	}

	// Collect index statements to add after table/field definitions
	var indexStatements []string

	for _, node := range b.input.nodes {
		statement := fmt.Sprintf("DEFINE TABLE %s SCHEMAFULL TYPE NORMAL PERMISSIONS FULL;", node.NameDatabase())
		statements = append(statements, statement)

		for _, f := range node.GetFields() {
			statements = append(statements, f.SchemaStatements(node.NameDatabase(), "")...)
		}

		// Build indexes for this table (handles both simple and composite)
		indexStatements = append(indexStatements, b.buildTableIndexStatements(node.NameDatabase(), node.GetFields(), node.Source.SoftDelete)...)

		// Index expires_at to keep expiry purge deletes efficient.
		if node.Source.Expiry {
			indexName := fmt.Sprintf(def.IndexPrefix+"%s_expires_at", node.NameDatabase())
			indexStatements = append(indexStatements,
				fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS expires_at CONCURRENTLY;", indexName, node.NameDatabase()))
		}

		statements = append(statements, "")
	}

	for _, edge := range b.input.edges {
		statement := fmt.Sprintf(
			"DEFINE TABLE %s SCHEMAFULL TYPE RELATION IN %s OUT %s ENFORCED PERMISSIONS FULL;",
			edge.NameDatabase(),
			edge.In.NameDatabase(),
			edge.Out.NameDatabase(), // TODO: can be OR'ed with "|"
		)
		statements = append(statements, statement)

		for _, f := range edge.GetFields() {
			statements = append(statements, f.SchemaStatements(edge.NameDatabase(), "")...)
		}

		// Build indexes for this table (handles both simple and composite)
		indexStatements = append(indexStatements, b.buildTableIndexStatements(edge.NameDatabase(), edge.GetFields(), edge.Source.SoftDelete)...)

		statements = append(statements, "")
	}

	// Sinks are write-only ingestion tables: records are accepted (firing
	// any dependent views/events) but discarded via DROP. Fields are still
	// defined so writes are validated and dependent view SELECTs typecheck.
	// No indexes are emitted — there are no rows to index.
	for _, sink := range b.input.sinks {
		statement := fmt.Sprintf("DEFINE TABLE %s DROP SCHEMAFULL TYPE NORMAL PERMISSIONS FULL;", sink.NameDatabase())
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

	var sb strings.Builder
	fmt.Fprintf(&sb, "DEFINE TABLE %s TYPE NORMAL AS SELECT %s FROM %s",
		view.NameDatabase(),
		strings.Join(def.Projections, ", "),
		sourceDB,
	)

	if def.Where != "" {
		fmt.Fprintf(&sb, " WHERE %s", def.Where)
	}

	if len(def.GroupBy) > 0 {
		fmt.Fprintf(&sb, " GROUP BY %s", strings.Join(def.GroupBy, ", "))
	}

	sb.WriteString(";")

	return sb.String(), nil
}

func buildAnalyzerStatement(analyzer parser.AnalyzerDef) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("DEFINE ANALYZER %s", analyzer.Name))

	if len(analyzer.Tokenizers) > 0 {
		parts = append(parts, fmt.Sprintf("TOKENIZERS %s", strings.Join(analyzer.Tokenizers, ", ")))
	}

	if len(analyzer.Filters) > 0 {
		var filterParts []string
		for _, filter := range analyzer.Filters {
			filterParts = append(filterParts, buildFilterString(filter))
		}
		parts = append(parts, fmt.Sprintf("FILTERS %s", strings.Join(filterParts, ", ")))
	}

	return strings.Join(parts, " ") + ";"
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
func (b *build) buildTableIndexStatements(tableName string, fields []field.Field, softDelete bool) []string {
	var statements []string

	if !b.noCountIndex {
		stmt := fmt.Sprintf("DEFINE INDEX "+def.IndexPrefix+"%s_count ON %s COUNT;", tableName, tableName)
		statements = append(statements, stmt)
	}

	// Collect composite unique index fields grouped by name
	compositeUnique := make(map[string][]string) // name -> []fieldPath

	// Process all fields (including nested)
	b.collectIndexes(tableName, "", fields, &statements, compositeUnique)

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

	return statements
}

// collectIndexes recursively collects index statements and composite unique fields.
func (b *build) collectIndexes(tableName, fieldPrefix string, fields []field.Field, statements *[]string, compositeUnique map[string][]string) {
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
				stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s FULLTEXT ANALYZER %s",
					indexName, tableName, fieldPath, searchDef.AnalyzerName)
				if searchDef.HasBM25 {
					stmt += fmt.Sprintf(" BM25(%g, %g)", searchDef.BM25K1, searchDef.BM25B)
				} else {
					stmt += " BM25"
				}
				if searchDef.Highlights {
					stmt += " HIGHLIGHTS"
				}
				if searchDef.Concurrently {
					stmt += " CONCURRENTLY"
				}
				stmt += ";"
				*statements = append(*statements, stmt)
			}
		}

		// Handle nested struct fields
		if nestedFields := f.NestedFields(); nestedFields != nil {
			b.collectIndexes(tableName, fieldPath, nestedFields, statements, compositeUnique)
		}
	}
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
