package codegen

const baseFileTmpl = `{{.ImportBlock}}

type {{.NameGo}}Repo interface {
	// Query returns a new query builder for the {{.NameGo}} model.

	Query() query.Builder[model.{{.NameGo}}]
{{- if not .HasComplexID}}
	// Create creates a new record for the {{.NameGo}} model.

	Create(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
{{- end}}
{{- if .HasComplexID}}
	// CreateWithID creates a new record with the given key for the {{.NameGo}} model.

	CreateWithID(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
{{- else}}
	// CreateWithID creates a new record with the given ID for the {{.NameGo}} model.

	CreateWithID(ctx context.Context, id string, {{.NameGoLower}} *model.{{.NameGo}}) error
{{- end}}
{{- if .HasComplexID}}
	// Read returns the record for the given key, if it exists.

	Read(ctx context.Context, key {{.KeyType}}) (*model.{{.NameGo}}, bool, error)
{{- else}}
	// Read returns the record for the given ID, if it exists.

	Read(ctx context.Context, id string) (*model.{{.NameGo}}, bool, error)
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

{{range .HookEvents}}	{{.Comment}}

	On{{.TimingUpper}}{{.Event}}(fn func(ctx context.Context, node *model.{{$.NameGo}}) error) func()
{{end -}}
}

// {{.NameGoLower}}RepoInfo holds the model-specific conversion functions for {{.NameGo}}.
var {{.NameGoLower}}RepoInfo = RepoInfo[model.{{.NameGo}}]{
	MarshalOne: func(node *model.{{.NameGo}}) any {
		return conv.From{{.NameGo}}Ptr(node)
	},
	UnmarshalOne: func(unmarshal func([]byte, any) error, data []byte) (*model.{{.NameGo}}, error) {
		var raw *conv.{{.NameGo}}
		if err := unmarshal(data, &raw); err != nil {
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
		c.{{.NameGoLower}}Repo = &{{.NameGoLower}}{repo: &repo[model.{{.NameGo}}, {{.KeyType}}]{
			db:   c.db,
			name: "{{.NameDB}}",
			info: {{.NameGoLower}}RepoInfo,
{{- if not .HasComplexID}}
			newID: {{.NewIDFunc}},
{{- end}}
			recordID: {{.RecordIDFunc}}}}
	}
	return c.{{.NameGoLower}}Repo
}

type {{.NameGoLower}} struct {
	*repo[model.{{.NameGo}}, {{.KeyType}}]
	mu           sync.RWMutex
{{- range .HookEvents}}
	{{.TimingLower}}{{.Event}} []{{$.NameGoLower}}Hook
{{- end}}
}

type {{.NameGoLower}}Hook struct {
	id uint64
	fn func(ctx context.Context, node *model.{{.NameGo}}) error
}

var {{.NameGoLower}}HookCounter atomic.Uint64
{{range .HookEvents}}
{{.Comment}}
func (r *{{$.NameGoLower}}) On{{.TimingUpper}}{{.Event}}(fn func(ctx context.Context, node *model.{{$.NameGo}}) error) func() {
	id := {{$.NameGoLower}}HookCounter.Add(1)
	r.mu.Lock()
	r.{{.TimingLower}}{{.Event}} = append(r.{{.TimingLower}}{{.Event}}, {{$.NameGoLower}}Hook{
		fn: fn,
		id: id,
	})
	r.mu.Unlock()
	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		for i, h := range r.{{.TimingLower}}{{.Event}} {
			if h.id == id {
				r.{{.TimingLower}}{{.Event}} = slices.Delete(r.{{.TimingLower}}{{.Event}}, i, i+1)
				return
			}
		}
	}
}
{{end -}}

// Query returns a new query builder for the {{.NameGo}} model.
func (r *{{.NameGoLower}}) Query() query.Builder[model.{{.NameGo}}] {
	return query.New{{.NameGo}}(r.db)
}
{{if not .HasComplexID}}
// Create creates a new record for the {{.NameGo}} model.
// The ID will be generated automatically as a ULID.
func (r *{{.NameGoLower}}) Create(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
	if {{.NameGoLower}} == nil {
		return errors.New("the passed node must not be nil")
	}
	if {{.NameGoLower}}.ID() != "" {
		return errors.New("given node already has an id")
	}
	{{call .BeforeHooks "Create"}}
	if err := r.create(ctx, {{.NameGoLower}}); err != nil {
		return err
	}
	{{call .AfterHooks "Create"}}
	return nil
}
{{end -}}
{{if .HasComplexID}}
// CreateWithID creates a new record for the {{.NameGo}} model using its embedded key.
// The node must have a non-zero ID set.
func (r *{{.NameGoLower}}) CreateWithID(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
	if {{.NameGoLower}} == nil {
		return errors.New("the passed node must not be nil")
	}
	{{call .IDCheck .NameGoLower "node must have a non-zero ID"}}
	{{call .BeforeHooks "Create"}}
	if err := r.createWithID(ctx, {{.NameGoLower}}.ID(), {{.NameGoLower}}); err != nil {
		return err
	}
	{{call .AfterHooks "Create"}}
	return nil
}
{{else}}
// CreateWithID creates a new record for the {{.NameGo}} model with the given id.
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
	{{call .BeforeHooks "Create"}}
	if err := r.createWithID(ctx, id, {{.NameGoLower}}); err != nil {
		return err
	}
	{{call .AfterHooks "Create"}}
	return nil
}
{{end -}}
{{if .HasComplexID}}
// Read returns the record for the given key, if it exists.
// The returned bool indicates whether the record was found or not.
func (r *{{.NameGoLower}}) Read(ctx context.Context, key {{.KeyType}}) (*model.{{.NameGo}}, bool, error) {
	if internal.CacheEnabled[model.{{.NameGo}}](ctx) {
		return nil, false, som.ErrCacheNotSupported
	}
	return r.read(ctx, r.recordID(key))
}
{{else}}
// Read returns the record for the given id, if it exists.
// The returned bool indicates whether the record was found or not.
// If caching is enabled via som.WithCache, it will be used.
func (r *{{.NameGoLower}}) Read(ctx context.Context, id string) (*model.{{.NameGo}}, bool, error) {
	if id == "" {
		return nil, false, som.ErrEmptyID
	}
	rid := r.recordID(id)
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
{{end}}
// Update updates the record for the given model.
func (r *{{.NameGoLower}}) Update(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
	if {{.NameGoLower}} == nil {
		return errors.New("the passed node must not be nil")
	}
	{{call .IDCheck .NameGoLower (printf "cannot update %s without existing record ID" .NameGo)}}
	{{call .BeforeHooks "Update"}}
	if err := r.update(ctx, {{.RecordIDFromVar}}, {{.NameGoLower}}); err != nil {
		return err
	}
	{{call .AfterHooks "Update"}}
	return nil
}

// Delete deletes the record for the given model.
func (r *{{.NameGoLower}}) Delete(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
	if {{.NameGoLower}} == nil {
		return errors.New("the passed node must not be nil")
	}
	{{call .IDCheck .NameGoLower (printf "cannot delete %s without existing record ID" .NameGo)}}
{{- if .SoftDelete}}
	if {{.NameGoLower}}.SoftDelete.IsDeleted() {
		return som.ErrAlreadyDeleted
	}
{{- end}}
	{{call .BeforeHooks "Delete"}}
{{- if and .SoftDelete .OptimisticLock}}
	version := {{.NameGoLower}}.Version()
	if err := r.delete(ctx, {{.RecordIDFromVar}}, {{.NameGoLower}}, true, &version); err != nil {
		return err
	}
{{- else}}
	if err := r.delete(ctx, {{.RecordIDFromVar}}, {{.NameGoLower}}, {{.SoftDelete}}, nil); err != nil {
		return err
	}
{{- end}}
	{{call .AfterHooks "Delete"}}
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
	{{call .IDCheck .NameGoLower (printf "cannot erase %s without existing record ID" .NameGo)}}
	return r.delete(ctx, {{.RecordIDFromVar}}, {{.NameGoLower}}, false, nil)
}

// Restore un-deletes a soft-deleted record.
// Sets deleted_at to NONE and refreshes the in-memory object.
func (r *{{.NameGoLower}}) Restore(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
	if {{.NameGoLower}} == nil {
		return errors.New("the passed node must not be nil")
	}
	{{call .IDCheck .NameGoLower (printf "cannot restore %s without existing record ID" .NameGo)}}
	if !{{.NameGoLower}}.SoftDelete.IsDeleted() {
		return errors.New("record is not deleted, cannot restore")
	}
{{- if .OptimisticLock}}
	query := "UPDATE $id SET deleted_at = NONE, __som_lock_version = $lock_version"
	vars := map[string]any{
		"id":           {{.RecordIDFromVar}},
		"lock_version": {{.NameGoLower}}.Version(),
	}
{{- else}}
	query := "UPDATE $id SET deleted_at = NONE"
	vars := map[string]any{"id": {{.RecordIDFromVar}}}
{{- end}}
	_, err := r.db.Query(ctx, query, vars)
{{- if .OptimisticLock}}
	if err != nil {
		if strings.Contains(err.Error(), "optimistic_lock_failed") {
			return som.ErrOptimisticLock
		}
		return fmt.Errorf("could not restore entity: %w", err)
	}
{{- else}}
	if err != nil {
		return fmt.Errorf("could not restore entity: %w", err)
	}
{{- end}}
	return r.refresh(ctx, {{.RecordIDFromVar}}, {{.NameGoLower}})
}
{{end -}}

// Refresh refreshes the given model with the remote data.
func (r *{{.NameGoLower}}) Refresh(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
	if {{.NameGoLower}} == nil {
		return errors.New("the passed node must not be nil")
	}
	{{call .IDCheck .NameGoLower (printf "cannot refresh %s without existing record ID" .NameGo)}}
	return r.refresh(ctx, {{.RecordIDFromVar}}, {{.NameGoLower}})
}
{{if not .HasComplexID}}
// Relate returns a new relate instance for the {{.NameGo}} model.
func (r *{{.NameGoLower}}) Relate() *relate.{{.NameGo}} {
	return relate.New{{.NameGo}}(r.db)
}
{{end}}`
