package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
)

func (b *build) buildNodeRepoFile(node *field.NodeTable) error {
	tmpl := `
		type {{.NameGo}}Repo interface {
			// Query returns a new query builder for the {{.NameGo}} model.

			Query() query.Builder[model.{{.NameGo}}]
			{{- if not .HasComplexID}}
			// Create creates a new record for the {{.NameGo}} model.

			Create(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			// Insert creates multiple records in a single operation.
			// Before- and after-create hooks are invoked for each node.

			Insert(ctx context.Context, nodes []*model.{{.NameGo}}) error
			// CreateWithID creates a new record with the given ID for the {{.NameGo}} model.

			CreateWithID(ctx context.Context, id string, {{.NameGoLower}} *model.{{.NameGo}}) error
			// Read returns the record for the given ID, if it exists.

			Read(ctx context.Context, id string) (*model.{{.NameGo}}, bool, error)
			{{- else}}
			// CreateWithID creates a new record with the given key for the {{.NameGo}} model.

			CreateWithID(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			// Read returns the record for the given key, if it exists.

			Read(ctx context.Context, key {{.KeyType}}) (*model.{{.NameGo}}, bool, error)
			{{- end}}
			// Update updates the record for the given {{.NameGo}} model.

			Update(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			// Delete deletes the record for the given {{.NameGo}} model.

			Delete(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			{{- if .SoftDelete}}
			// Erase permanently deletes the record from the database.

			Erase(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			// Restore un-deletes a soft-deleted record.

			Restore(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			{{- end}}
			// Refresh refreshes the given model with the current database state.

			Refresh(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			{{- if not .HasComplexID}}
			// Relate returns a new relate builder for the {{.NameGo}} model.

			Relate() *relate.{{.NameGo}}
			{{- end}}
			// Index returns a new index instance for the {{.NameGo}} model.

			Index() *index.{{.NameGo}}
			{{- if .HasChangefeed}}
			// Changes returns a new changes query builder for the {{.NameGo}} model.
			// This is only available for models with changefeed enabled.

			Changes() query.ChangesBuilder[model.{{.NameGo}}, conv.{{.NameGo}}]
			{{- end}}

			{{range .Hooks}}{{.Comment}}

			{{.Name}}(fn func(ctx context.Context, node *model.{{$.NameGo}}) error) func()
			{{end -}}
		}

		// {{.NameGoLower}}RepoInfo holds the model-specific conversion functions for {{.NameGo}}.
		var {{.NameGoLower}}RepoInfo = RepoInfo[model.{{.NameGo}}]{
			CreateNew: func(ctx context.Context, db *dbConn, target string, data any) (*model.{{.NameGo}}, error) {
				raw, err := dbCreateNew[conv.{{.NameGo}}](ctx, db, target, data)
				if err != nil {
					return nil, err
				}
				return conv.To{{.NameGo}}Ptr(raw), nil
			},
			CreateOne: func(ctx context.Context, db *dbConn, id models.RecordID, data any) (*model.{{.NameGo}}, error) {
				raw, err := dbCreate[conv.{{.NameGo}}](ctx, db, id, data)
				if err != nil {
					return nil, err
				}
				return conv.To{{.NameGo}}Ptr(raw), nil
			},
			InsertAll: func(ctx context.Context, db *dbConn, stmt string, vars map[string]any) ([]*model.{{.NameGo}}, error) {
				raw, err := dbInsert[conv.{{.NameGo}}](ctx, db, stmt, vars)
				if err != nil {
					return nil, err
				}
				results := make([]*model.{{.NameGo}}, len(raw))
				for i, r := range raw {
					results[i] = conv.To{{.NameGo}}Ptr(r)
				}
				return results, nil
			},
			MarshalOne: func(node *model.{{.NameGo}}) any {
				return conv.From{{.NameGo}}Ptr(node)
			},
			QueryOne: func(ctx context.Context, db *dbConn, stmt string, vars map[string]any) (*model.{{.NameGo}}, error) {
				raw, err := dbQueryOne[conv.{{.NameGo}}](ctx, db, stmt, vars)
				if err != nil {
					return nil, err
				}
				if raw == nil {
					return nil, nil
				}
				return conv.To{{.NameGo}}Ptr(raw), nil
			},
			ReadOne: func(ctx context.Context, db *dbConn, id *models.RecordID) (*model.{{.NameGo}}, error) {
				raw, err := dbSelect[conv.{{.NameGo}}](ctx, db, id)
				if err != nil {
					return nil, err
				}
				if raw == nil {
					return nil, nil
				}
				return conv.To{{.NameGo}}Ptr(raw), nil
			},
			UpdateOne: func(ctx context.Context, db *dbConn, id *models.RecordID, data any) (*model.{{.NameGo}}, error) {
				raw, err := dbUpdate[conv.{{.NameGo}}](ctx, db, id, data)
				if err != nil {
					return nil, err
				}
				return conv.To{{.NameGo}}Ptr(raw), nil
			},
		}

		// {{.NameGo}}Repo returns the repository instance for the {{.NameGo}} model.
		// The instance is cached as a singleton on the client.
		func (c *ClientImpl) {{.NameGo}}Repo() {{.NameGo}}Repo {
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.{{.NameGoLower}}Repo == nil {
				c.{{.NameGoLower}}Repo = &{{.NameGoLower}}{repo: {{.RepoLiteral}}}
			}
			return c.{{.NameGoLower}}Repo
		}

		type {{.NameGoLower}} struct {
			*repo[model.{{.NameGo}}, {{.KeyType}}]
		}

		// Query returns a new query builder for the {{.NameGo}} model.
		func (r *{{.NameGoLower}}) Query() query.Builder[model.{{.NameGo}}] {
			return query.New{{.NameGo}}(r.db)
		}
		{{if not .HasComplexID}}
		// Create creates a new record for the {{.NameGo}} model.
		// The ID will be generated automatically as a ULID.
		// Before- and after-create hooks are invoked.
		func (r *{{.NameGoLower}}) Create(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			if {{.NameGoLower}}.ID() != "" {
				return errors.New("given node already has an id")
			}
			{{call .RunHooks "beforeCreate"}}
			if err := r.create(ctx, {{.NameGoLower}}); err != nil {
				return err
			}
			{{call .RunHooks "afterCreate"}}
			return nil
		}

		// CreateWithID creates a new record for the {{.NameGo}} model with the given id.
		// Before- and after-create hooks are invoked.
		func (r *{{.NameGoLower}}) CreateWithID(ctx context.Context, id string, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			if id == "" {
				return som.ErrEmptyID
			}
			if {{.NameGoLower}}.ID() != "" {
				return errors.New("given node already has an id")
			}
			{{call .RunHooks "beforeCreate"}}
			if err := r.createWithID(ctx, id, {{.NameGoLower}}); err != nil {
				return err
			}
			{{call .RunHooks "afterCreate"}}
			return nil
		}

		// Insert creates multiple records in a single operation.
		// Before- and after-create hooks are invoked for each node.
		func (r *{{.NameGoLower}}) Insert(ctx context.Context, nodes []*model.{{.NameGo}}) error {
			if len(nodes) == 0 {
				return nil
			}
			for _, n := range nodes {
				if n == nil {
					return errors.New("slice contains nil node")
				}
				if n.ID() != "" {
					return errors.New("node already has an id")
				}
			}
			if err := r.runHooksAll(ctx, beforeCreate, nodes); err != nil {
				return err
			}
			if err := r.insert(ctx, nodes); err != nil {
				return err
			}
			if err := r.runHooksAll(ctx, afterCreate, nodes); err != nil {
				return err
			}
			return nil
		}

		// Read returns the record for the given id, if it exists.
		// The returned bool indicates whether the record was found or not.
		// If caching is enabled via som.WithCache, it will be used.
		func (r *{{.NameGoLower}}) Read(ctx context.Context, id string) (*model.{{.NameGo}}, bool, error) {
			if id == "" {
				return nil, false, som.ErrEmptyID
			}
			rid := r.recordID(id)
			if internal.TxActive(ctx) {
				return r.read(ctx, rid)
			}
			if !internal.CacheEnabled[model.{{.NameGo}}](ctx) {
				return r.read(ctx, rid)
			}
			idFunc := func(n *model.{{.NameGo}}) string {
				return string(n.ID())
			}
			queryAll := func(ctx context.Context) ([]*model.{{.NameGo}}, error) {
				return r.Query().All(ctx)
			}
			countAll := func(ctx context.Context) (int, error) {
				return r.Query().Count(ctx)
			}
			cache, err := getOrCreateCache[model.{{.NameGo}}](ctx, idFunc, queryAll, countAll)
			if err != nil {
				return nil, false, err
			}
			var refreshFuncs *eagerRefreshFuncs[model.{{.NameGo}}]
			if cache != nil && cache.isEager() {
				refreshFuncs = &eagerRefreshFuncs[model.{{.NameGo}}]{cacheID: internal.GetCacheKey[model.{{.NameGo}}](ctx), queryAll: queryAll, countAll: countAll, idFunc: idFunc}
			}
			return r.readWithCache(ctx, id, rid, cache, refreshFuncs)
		}
		{{else}}
		// CreateWithID creates a new record for the {{.NameGo}} model using its embedded key.
		// The node must have a non-zero ID set.
		// Before- and after-create hooks are invoked.
		func (r *{{.NameGoLower}}) CreateWithID(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			{{call .IDCheck "node must have a non-zero ID"}}
			{{call .RunHooks "beforeCreate"}}
			if err := r.createWithID(ctx, {{.NameGoLower}}.ID(), {{.NameGoLower}}); err != nil {
				return err
			}
			{{call .RunHooks "afterCreate"}}
			return nil
		}

		// Read returns the record for the given key, if it exists.
		// The returned bool indicates whether the record was found or not.
		func (r *{{.NameGoLower}}) Read(ctx context.Context, key {{.KeyType}}) (*model.{{.NameGo}}, bool, error) {
			if internal.CacheEnabled[model.{{.NameGo}}](ctx) {
				return nil, false, som.ErrCacheNotSupported
			}
			return r.read(ctx, r.recordID(key))
		}
		{{end}}
		// Update updates the record for the given model.
		// Before- and after-update hooks are invoked.
		func (r *{{.NameGoLower}}) Update(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			{{call .IDCheck (printf "cannot update %s without existing record ID" .NameGo)}}
			{{call .RunHooks "beforeUpdate"}}
			if err := r.update(ctx, {{.RecordIDFromNode}}, {{.NameGoLower}}); err != nil {
				return err
			}
			{{call .RunHooks "afterUpdate"}}
			return nil
		}

		// Delete deletes the record for the given model.
		// Before- and after-delete hooks are invoked.
		func (r *{{.NameGoLower}}) Delete(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			{{call .IDCheck (printf "cannot delete %s without existing record ID" .NameGo)}}
			{{- if .SoftDelete}}
			if {{.NameGoLower}}.SoftDelete.IsDeleted() {
				return som.ErrAlreadyDeleted
			}
			{{- end}}
			{{call .RunHooks "beforeDelete"}}
			{{- if and .SoftDelete .OptimisticLock}}
			version := {{.NameGoLower}}.Version()
			if err := r.delete(ctx, {{.RecordIDFromNode}}, {{.NameGoLower}}, true, &version); err != nil {
				return err
			}
			{{- else}}
			if err := r.delete(ctx, {{.RecordIDFromNode}}, {{.NameGoLower}}, {{.SoftDelete}}, nil); err != nil {
				return err
			}
			{{- end}}
			{{call .RunHooks "afterDelete"}}
			return nil
		}
		{{if .SoftDelete}}
		// Erase permanently deletes the record from the database.
		// This performs a hard delete and cannot be undone.
		// Use this to permanently remove soft-deleted records.
		func (r *{{.NameGoLower}}) Erase(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			{{call .IDCheck (printf "cannot erase %s without existing record ID" .NameGo)}}
			return r.delete(ctx, {{.RecordIDFromNode}}, {{.NameGoLower}}, false, nil)
		}

		// Restore un-deletes a soft-deleted record.
		// Sets deleted_at to NONE and refreshes the in-memory object.
		func (r *{{.NameGoLower}}) Restore(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			{{call .IDCheck (printf "cannot restore %s without existing record ID" .NameGo)}}
			if !{{.NameGoLower}}.SoftDelete.IsDeleted() {
				return errors.New("record is not deleted, cannot restore")
			}
			{{- if .OptimisticLock}}
			query := "UPDATE $id SET deleted_at = NONE, __som_lock_version = $lock_version"
			vars := map[string]any{
				"id":           {{.RecordIDFromNode}},
				"lock_version": {{.NameGoLower}}.Version(),
			}
			{{- else}}
			query := "UPDATE $id SET deleted_at = NONE"
			vars := map[string]any{"id": {{.RecordIDFromNode}}}
			{{- end}}
			result, err := r.info.QueryOne(ctx, r.db, query, vars)
			if err != nil {
				{{- if .OptimisticLock}}
				if containsError(err, "optimistic_lock_failed") {
					return fmt.Errorf("%w: %w", som.ErrOptimisticLock, err)
				}
				{{- end}}
				return fmt.Errorf("could not restore entity: %w", err)
			}
			if result == nil {
				return som.ErrNotFound
			}
			*{{.NameGoLower}} = *result
			return nil
		}
		{{end}}
		// Refresh refreshes the given model with the remote data.
		func (r *{{.NameGoLower}}) Refresh(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed node must not be nil")
			}
			{{call .IDCheck (printf "cannot refresh %s without existing record ID" .NameGo)}}
			return r.refresh(ctx, {{.RecordIDFromNode}}, {{.NameGoLower}})
		}
		{{if not .HasComplexID}}
		// Relate returns a new relate instance for the {{.NameGo}} model.
		func (r *{{.NameGoLower}}) Relate() *relate.{{.NameGo}} {
			return relate.New{{.NameGo}}(r.db)
		}
		{{end}}
		// Index returns a new index instance for the {{.NameGo}} model.
		func (r *{{.NameGoLower}}) Index() *index.{{.NameGo}} {
			return index.New{{.NameGo}}(r.db)
		}
		{{if .HasChangefeed}}
		// Changes returns a new changes query builder for the {{.NameGo}} model.
		// This is only available for models with changefeed enabled.
		func (r *{{.NameGoLower}}) Changes() query.ChangesBuilder[model.{{.NameGo}}, conv.{{.NameGo}}] {
			return query.New{{.NameGo}}Changes(r.db)
		}
		{{end -}}
	`

	data := map[string]any{
		"NameGo":           node.NameGo(),
		"NameGoLower":      node.NameGoLower(),
		"NameDB":           node.NameDatabase(),
		"KeyType":          b.keyType(node),
		"HasComplexID":     node.HasComplexID(),
		"HasChangefeed":    node.HasChangefeed(),
		"SoftDelete":       node.Source.SoftDelete,
		"OptimisticLock":   node.Source.OptimisticLock,
		"RepoLiteral":      b.repoLiteral(node),
		"RecordIDFromNode": b.recordIDFromNode(node),
		"Hooks":            repoHooks(),

		"IDCheck":  func(errMsg string) string { return b.idEmptyCheck(node, errMsg) },
		"RunHooks": func(kind string) string { return runHooksCall(node, kind) },
	}

	return renderGoFileWithImports(
		b.fs.Writer(filepath.Join(def.PkgRepo, node.FileName())),
		def.PkgRepo, "repoNode", tmpl, data,
		[]goImport{
			{Path: "context"},
			{Path: "errors"},
			{Path: "fmt"},
			{Alias: "models", Path: def.PkgModels},
			{Alias: "som", Path: b.relativePkgPath()},
			{Alias: "conv", Path: b.relativePkgPath(def.PkgConv)},
			{Alias: "index", Path: b.relativePkgPath(def.PkgIndex)},
			{Alias: "internal", Path: b.relativePkgPath(def.PkgInternal)},
			{Alias: "types", Path: b.relativePkgPath(def.PkgTypes)},
			{Alias: "query", Path: b.relativePkgPath(def.PkgQuery)},
			{Alias: "relate", Path: b.relativePkgPath(def.PkgRelate)},
			{Alias: "model", Path: b.input.sourcePkgPath},
		},
	)
}

// repoHook describes one of the hook registration methods of a repository.
type repoHook struct {
	Name    string
	Comment string
}

// repoHooks returns the hook registration methods of a node repository,
// in the order they are generated into the repository interface.
func repoHooks() []repoHook {
	var hooks []repoHook

	for _, event := range []string{"Create", "Update", "Delete"} {
		for _, timing := range []string{"Before", "After"} {
			name := "On" + timing + event
			verb := strings.ToLower(event)

			effect := "runs after a record has been " + verb + "d.\n" +
				"If the hook returns an error, the error is returned to the caller."

			if timing == "Before" {
				effect = "runs before a record is " + verb + "d.\n" +
					"If the hook returns an error, the " + verb + " operation is aborted."
			}

			hooks = append(hooks, repoHook{
				Name: name,
				Comment: formatGoComment(name + " registers a hook that " + effect + "\n" +
					"Returns a function that, when called, removes this hook.\n" +
					"\n" +
					"Note: Hooks are local to this application instance and are not\n" +
					"distributed across multiple instances of the application."),
			})
		}
	}

	return hooks
}

// runHooksCall emits a call to the generic repo.runHooks for the given hook
// kind against the node under operation, returning on error.
func runHooksCall(node *field.NodeTable, kind string) string {
	return renderCode(jen.If(
		jen.Err().Op(":=").Id("r").Dot("runHooks").Call(
			jen.Id("ctx"), jen.Id(kind), jen.Id(node.NameGoLower()),
		),
		jen.Err().Op("!=").Nil(),
	).Block(jen.Return(jen.Err())))
}

// keyType returns the type of the record key, which is the complex ID struct
// for nodes with a complex ID and a plain string for all other nodes.
func (b *build) keyType(node *field.NodeTable) string {
	if node.HasComplexID() {
		return "model." + node.Source.ComplexID.StructName
	}

	return "string"
}

// repoLiteral returns the composite literal for the generic repo the node
// repository is built upon.
func (b *build) repoLiteral(node *field.NodeTable) string {
	values := []jen.Code{
		jen.Line().Id("db").Op(":").Id("c").Dot("db"),
		jen.Line().Id("name").Op(":").Lit(node.NameDatabase()),
		jen.Line().Id("info").Op(":").Id(node.NameGoLower() + "RepoInfo"),
	}

	if !node.HasComplexID() {
		values = append(values, jen.Line().Id("autoID").Op(":").True())
	}

	values = append(values,
		jen.Line().Id("recordID").Op(":").Add(b.recordIDFunc(node)),
		jen.Line(),
	)

	return renderCode(jen.Op("&").Id("repo").
		Types(jen.Id("model."+node.NameGo()), jen.Id(b.keyType(node))).
		Values(values...))
}

// recordIDFunc returns the function literal that converts a record key into
// the record ID of the database.
func (b *build) recordIDFunc(node *field.NodeTable) jen.Code {
	recordID := jen.Op("*").Qual(def.PkgModels, "RecordID")

	if node.HasComplexID() {
		return jen.Func().Params(jen.Id("key").Id(b.keyType(node))).Add(recordID).Block(
			jen.Id("rid").Op(":=").Qual(def.PkgModels, "NewRecordID").Call(
				jen.Lit(node.NameDatabase()),
				b.recordIDValue(node, "key"),
			),
			jen.Return(jen.Op("&").Id("rid")),
		)
	}

	parseFunc := "parseStringID"
	if node.Source.IDType == parser.IDTypeUUID {
		parseFunc = "parseUUID"
	}

	return jen.Func().Params(jen.Id("id").String()).Add(recordID).Block(
		jen.Id("rid").Op(":=").Qual(def.PkgModels, "NewRecordID").Call(
			jen.Lit(node.NameDatabase()),
			jen.Id(parseFunc).Call(jen.Id("id")),
		),
		jen.Return(jen.Op("&").Id("rid")),
	)
}

// recordIDFromNode returns the expression building the record ID from the ID
// of the node under operation.
func (b *build) recordIDFromNode(node *field.NodeTable) string {
	if node.HasComplexID() {
		return fmt.Sprintf("r.recordID(%s.ID())", node.NameGoLower())
	}

	return fmt.Sprintf("r.recordID(string(%s.ID()))", node.NameGoLower())
}

// idEmptyCheck emits the guard clauses ensuring that the ID of the node under
// operation is set. For complex IDs referencing other nodes, the ID of each
// referenced node is checked instead.
func (b *build) idEmptyCheck(node *field.NodeTable, errMsg string) string {
	id := jen.Id(node.NameGoLower()).Dot("ID").Call()

	if !node.HasComplexID() {
		return renderCode(emptyStringCheck(id, errMsg))
	}

	cid := node.Source.ComplexID

	if !cid.HasNodeRef() {
		return renderCode(jen.Var().Id("zeroKey").Id(b.keyType(node)).Line().
			Add(zeroValueCheck(id, "zeroKey", errMsg)))
	}

	var checks []jen.Code

	for _, sf := range cid.Fields {
		fieldNode, ok := sf.Field.(*parser.FieldNode)
		if !ok {
			continue
		}

		refNode := b.input.findNodeByName(fieldNode.Node)
		if refNode == nil {
			continue
		}

		refID := jen.Add(id.Clone()).Dot(sf.Name).Dot("ID").Call()
		refErrMsg := sf.Name + ".ID must not be empty"

		switch {
		case !refNode.HasComplexID():
			checks = append(checks, emptyStringCheck(refID, refErrMsg))

		case !refNode.Source.ComplexID.HasNodeRef():
			zeroVar := "zero" + sf.Name + "Key"

			checks = append(checks, jen.Var().Id(zeroVar).Id(b.keyType(refNode)).Line().
				Add(zeroValueCheck(refID, zeroVar, refErrMsg)))
		}
	}

	return renderCode(joinStatements(checks))
}

// joinStatements combines the given statements into a single code fragment,
// one statement per line.
func joinStatements(codes []jen.Code) jen.Code {
	stmt := jen.Null()

	for i, code := range codes {
		if i > 0 {
			stmt.Line()
		}

		stmt.Add(code)
	}

	return stmt
}

func emptyStringCheck(expr jen.Code, errMsg string) jen.Code {
	return zeroValueCheck(expr, `""`, errMsg)
}

func zeroValueCheck(expr jen.Code, zeroValue, errMsg string) jen.Code {
	return jen.If(jen.Add(expr).Op("==").Id(zeroValue)).Block(
		jen.Return(jen.Qual("errors", "New").Call(jen.Lit(errMsg))),
	)
}

func (b *build) recordIDValue(node *field.NodeTable, keyVar string) jen.Code {
	cid := node.Source.ComplexID

	if cid.Kind == parser.IDTypeArray {
		var elems []jen.Code
		for _, sf := range cid.Fields {
			elems = append(elems, b.fieldValue(sf, keyVar))
		}
		return jen.Index().Any().Values(elems...)
	}

	dict := jen.Dict{}
	for _, sf := range cid.Fields {
		dict[jen.Lit(sf.DBName)] = b.fieldValue(sf, keyVar)
	}
	return jen.Map(jen.String()).Any().Values(dict)
}

func (b *build) fieldValue(sf parser.ComplexIDField, keyVar string) jen.Code {
	accessor := jen.Id(keyVar).Dot(sf.Name)
	return fieldValueFrom(b.input, b.basePkg(), sf, accessor)
}
