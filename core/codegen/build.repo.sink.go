package codegen

import (
	"path/filepath"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
)

// buildSinkRepoFile generates the write-only repository for a sink. Sinks
// expose only Create and Insert — the written rows are discarded by the
// DROP table, so nothing is read back, queried, updated or deleted.
func (b *build) buildSinkRepoFile(sink *field.SinkTable) error {
	tmpl := `
		type {{.NameGo}}Repo interface {
			// Create writes a new {{.NameGo}} record. The record is discarded
			// immediately after write (DROP table); nothing is returned.
			Create(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error
			// Insert writes multiple {{.NameGo}} records in a single operation.
			// The records are discarded immediately after write (DROP table).
			Insert(ctx context.Context, {{.NameGoLower}}s []*model.{{.NameGo}}) error
		}

		// {{.NameGo}}Repo returns the repository instance for the {{.NameGo}} sink.
		// The instance is cached as a singleton on the client.
		func (c *ClientImpl) {{.NameGo}}Repo() {{.NameGo}}Repo {
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.{{.NameGoLower}}Repo == nil {
				c.{{.NameGoLower}}Repo = &{{.NameGoLower}}{db: c.db}
			}
			return c.{{.NameGoLower}}Repo
		}

		type {{.NameGoLower}} struct {
			db *dbConn
		}

		// Create writes a new {{.NameGo}} record; the row is discarded after write.
		func (r *{{.NameGoLower}}) Create(ctx context.Context, {{.NameGoLower}} *model.{{.NameGo}}) error {
			if {{.NameGoLower}} == nil {
				return errors.New("the passed record must not be nil")
			}
			return dbInsertVoid(ctx, r.db, "{{.NameDB}}", []any{conv.From{{.NameGo}}Ptr({{.NameGoLower}})})
		}

		// Insert writes multiple {{.NameGo}} records; the rows are discarded after write.
		func (r *{{.NameGoLower}}) Insert(ctx context.Context, {{.NameGoLower}}s []*model.{{.NameGo}}) error {
			if len({{.NameGoLower}}s) == 0 {
				return nil
			}
			data := make([]any, len({{.NameGoLower}}s))
			for i, s := range {{.NameGoLower}}s {
				if s == nil {
					return errors.New("slice contains nil record")
				}
				data[i] = conv.From{{.NameGo}}Ptr(s)
			}
			return dbInsertVoid(ctx, r.db, "{{.NameDB}}", data)
		}
	`

	data := map[string]any{
		"NameGo":      sink.NameGo(),
		"NameGoLower": sink.NameGoLower(),
		"NameDB":      sink.NameDatabase(),
	}

	file := newGoFile(def.PkgRepo,
		goImport{Path: "context"},
		goImport{Path: "errors"},
		goImport{Alias: "conv", Path: b.relativePkgPath(def.PkgConv)},
		goImport{Alias: "model", Path: b.input.sourcePkgPath},
	)

	return file.render(
		b.fs.Writer(filepath.Join(def.PkgRepo, sink.FileName())),
		"repoSink", tmpl, data,
	)
}
