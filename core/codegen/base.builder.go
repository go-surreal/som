package codegen

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

const (
	filenameInterfaces = "som.interfaces.go"
	filenameSchema     = "tables.surql"
)

type build struct {
	input  *input
	fs     *fs.FS
	outPkg string
}

func BuildStatic(fs *fs.FS, outPkg string) error {
	tmpl := &embed.Template{
		GenerateOutPath: outPkg,
	}

	files, err := embed.Read(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		fs.Write(file.Path, file.Content)
	}

	return nil
}

func Build(source *parser.Output, fs *fs.FS, outPkg string) error {
	in, err := newInput(source, outPkg)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:  in,
		fs:     fs,
		outPkg: outPkg,
	}

	return builder.build()
}

func (b *build) build() error {
	if err := b.buildInterfaceFile(); err != nil {
		return err
	}

	if err := b.buildSchemaFile(); err != nil {
		return err
	}

	for _, node := range b.input.nodes {
		if err := b.buildBaseFile(node); err != nil {
			return err
		}
	}

	builders := []builder{
		b.newQueryBuilder(),
		b.newFilterBuilder(),
		b.newSortBuilder(),
		b.newFetchBuilder(),
		b.newConvBuilder(),
		b.newRelateBuilder(),
	}

	for _, builder := range builders {
		if err := builder.build(); err != nil {
			return err
		}
	}

	return nil
}

func (b *build) buildInterfaceFile() error {
	f := jen.NewFile(def.PkgRepo)

	f.PackageComment(string(embed.CodegenComment))

	f.Type().Id("Client").InterfaceFunc(func(g *jen.Group) {
		for _, node := range b.input.nodes {
			g.Id(node.NameGo() + "Repo").Call().Id(node.NameGo() + "Repo")
		}

		g.Id("ApplySchema").Call(jen.Id("ctx").Qual("context", "Context")).Error()
		g.Id("Close").Call()
	})

	if err := f.Render(b.fs.Writer(filepath.Join(def.PkgRepo, filenameInterfaces))); err != nil {
		return err
	}

	return nil
}

func (b *build) buildSchemaFile() error {
	statements := []string{string(embed.CodegenComment), ""}

	// Generate DEFINE ANALYZER statements first
	if b.input.config != nil {
		for _, analyzer := range b.input.config.Analyzers {
			statements = append(statements, buildAnalyzerStatement(analyzer))
		}
		if len(b.input.config.Analyzers) > 0 {
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
		indexStatements = append(indexStatements, b.buildTableIndexStatements(node.NameDatabase(), node.GetFields())...)

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
		indexStatements = append(indexStatements, b.buildTableIndexStatements(edge.NameDatabase(), edge.GetFields())...)

		statements = append(statements, "")
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
// simple indexes and composite unique indexes (fields grouped by UniqueName).
func (b *build) buildTableIndexStatements(tableName string, fields []field.Field) []string {
	var statements []string

	// Collect composite unique index fields grouped by UniqueName
	compositeUnique := make(map[string][]string) // UniqueName -> []fieldPath

	// Process all fields (including nested)
	b.collectIndexes(tableName, "", fields, &statements, compositeUnique)

	// Generate composite unique index statements
	for uniqueName, fieldPaths := range compositeUnique {
		// Index name format: __som_<table>_unique_<name>
		indexName := fmt.Sprintf("__som_%s_unique_%s", tableName, uniqueName)
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

		indexInfo := f.IndexInfo()
		searchInfo := f.SearchInfo()

		if indexInfo != nil {
			if indexInfo.Unique && indexInfo.UniqueName != "" {
				// Composite unique index - collect field for later
				compositeUnique[indexInfo.UniqueName] = append(compositeUnique[indexInfo.UniqueName], fieldPath)
			} else if indexInfo.Unique {
				// Simple unique index on single field
				// Index name format: __som_<table>_unique_<field>
				indexName := indexInfo.Name
				if indexName == "" {
					indexName = fmt.Sprintf("__som_%s_unique_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				}
				stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s UNIQUE;", indexName, tableName, fieldPath)
				*statements = append(*statements, stmt)
			} else {
				// Regular (non-unique) index
				// Index name format: __som_<table>_index_<field>
				indexName := indexInfo.Name
				if indexName == "" {
					indexName = fmt.Sprintf("__som_%s_index_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				}
				stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s CONCURRENTLY;", indexName, tableName, fieldPath)
				*statements = append(*statements, stmt)
			}
		}

		if searchInfo != nil && searchInfo.ConfigName != "" {
			// Look up the search config to get analyzer and options
			searchDef := b.findSearchConfig(searchInfo.ConfigName)
			if searchDef != nil {
				// Index name format: __som_<table>_search_<field>
				indexName := fmt.Sprintf("__som_%s_search_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				stmt := fmt.Sprintf("DEFINE INDEX %s ON %s FIELDS %s SEARCH ANALYZER %s",
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
	if b.input.config == nil {
		return nil
	}
	for i := range b.input.config.Searches {
		if b.input.config.Searches[i].Name == name {
			return &b.input.config.Searches[i]
		}
	}
	return nil
}

func (b *build) buildBaseFile(node *field.NodeTable) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(def.PkgRepo)

	f.PackageComment(string(embed.CodegenComment))

	//
	// type {NodeName}Repo interface {...}
	//
	f.Type().Id(node.NameGo()+"Repo").Interface(
		jen.Id("Query").Call().Qual(pkgQuery, "Builder").
			Types(b.input.SourceQual(node.NameGo()), jen.Qual(b.subPkg(def.PkgConv), node.NameGo())),

		jen.Id("Create").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("CreateWithID").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Read").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").Op("*").Qual(b.subPkg(""), "ID"),
		).Parens(jen.List(
			jen.Op("*").Add(b.input.SourceQual(node.NameGo())),
			jen.Bool(),
			jen.Error(),
		)),

		jen.Id("Update").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Delete").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Refresh").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Relate").Call().Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()),
	)

	f.Line().
		Add(comment(`
` + node.NameGo() + `Repo returns a new repository instance for the ` + node.NameGo() + ` model.
		`)).
		Func().Params(jen.Id("c").Op("*").Id("ClientImpl")).
		Id(node.NameGo() + "Repo").Params().Id(node.NameGo() + "Repo").
		Block(
			jen.Return(
				jen.Op("&").Id(node.NameGoLower()).Values(
					jen.Id("repo").Op(":").Op("&").Id("repo").
						Types(
							b.input.SourceQual(node.NameGo()),
							jen.Id("conv."+node.NameGo()),
						).
						Values(
							jen.Add(
								jen.Line(),
								jen.Id("db").Op(":").Id("c").Dot("db"),
							),
							jen.Add(
								jen.Line(),
								jen.Id("name").Op(":").Lit(node.NameDatabase()),
							),
							jen.Add(
								jen.Line(),
								jen.Id("convTo").Op(":").Qual(pkgConv, "To"+node.NameGo()+"Ptr"),
							),
							jen.Add(
								jen.Line(),
								jen.Id("convFrom").Op(":").Qual(pkgConv, "From"+node.NameGo()+"Ptr"),
							),
						),
				),
			),
		)

	f.Line()
	f.Type().Id(node.NameGoLower()).Struct(
		jen.Op("*").Id("repo").Types(
			b.input.SourceQual(node.NameGo()),
			jen.Id("conv."+node.NameGo()),
		),
	)

	f.Line().
		Add(comment(`
Query returns a new query builder for the `+node.NameGo()+` model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Query").Params().
		Qual(pkgQuery, "Builder").
		Types(
			b.input.SourceQual(node.NameGo()),
			jen.Qual(b.subPkg(def.PkgConv), node.NameGo()),
		).
		Block(
			jen.Return(jen.Qual(pkgQuery, "New"+node.NameGo()).Call(
				jen.Id("r").Dot("db"),
			)),
		)

	f.Line().
		Add(comment(`
Create creates a new record for the `+node.NameGo()+` model.
The ID will be generated automatically as a ULID.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))),
				),

			jen.Return(
				jen.Id("r").Dot("create").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
CreateWithID creates a new record for the `+node.NameGo()+` model with the given id.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("CreateWithID").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(), // TODO: name clash if node/model is named "id"!
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))),
				),

			jen.Return(
				jen.Id("r").Dot("createWithID").Call(
					jen.Id("ctx"),
					jen.Id("id"),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Read returns the record for the given id, if it exists.
The returned bool indicates whether the record was found or not.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").Op("*").Qual(b.subPkg(""), "ID"),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Bool(), jen.Error()).
		Block(
			jen.Return(
				jen.Id("r").Dot("read").Call(
					jen.Id("ctx"),
					jen.Id("id"),
				),
			),
		)

	f.Line().
		Add(comment(`
Update updates the record for the given model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Update").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot update "+node.NameGo()+" without existing record ID"))),
				),

			jen.Return(
				jen.Id("r").Dot("update").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Delete deletes the record for the given model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Delete").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.Return(
				jen.Id("r").Dot("delete").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Refresh refreshes the given model with the remote data.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Refresh").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot refresh "+node.NameGo()+" without existing record ID"))),
				),

			jen.Return(
				jen.Id("r").Dot("refresh").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Relate returns a new relate instance for the `+node.NameGo()+` model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Relate").Params().
		Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()).
		Block(
			jen.Return(
				jen.Qual(b.subPkg(def.PkgRelate), "New"+node.NameGo()).Call(
					jen.Id("r").Dot("db"),
				),
			),
		)

	if err := f.Render(b.fs.Writer(filepath.Join(def.PkgRepo, node.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *build) newQueryBuilder() builder {
	return newQueryBuilder(b.input, b.fs, b.basePkg(), def.PkgQuery)
}

func (b *build) newFilterBuilder() builder {
	return newFilterBuilder(b.input, b.fs, b.basePkg(), def.PkgFilter)
}

func (b *build) newSortBuilder() builder {
	return newSortBuilder(b.input, b.fs, b.basePkg(), def.PkgSort)
}

func (b *build) newFetchBuilder() builder {
	return newFetchBuilder(b.input, b.fs, b.basePkg(), def.PkgFetch)
}

func (b *build) newConvBuilder() builder {
	return newConvBuilder(b.input, b.fs, b.basePkg(), def.PkgConv)
}

func (b *build) newRelateBuilder() builder {
	return newRelateBuilder(b.input, b.fs, b.basePkg(), def.PkgRelate)
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) subPkg(pkg string) string {
	return path.Join(b.basePkg(), pkg)
}

//
// -- HELPER
//

func comment(text string) jen.Code {
	var code jen.Statement

	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		code.Comment(line).Line()
	}

	return &code
}
