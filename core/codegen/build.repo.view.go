package codegen

import (
	"path/filepath"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
)

// buildViewRepoFile generates the read-only repository for a view. Views
// expose only Query() — no create/update/delete/relate/index operations.
func (b *build) buildViewRepoFile(view *field.ViewTable) error {
	tmpl := `
		type {{.NameGo}}Repo interface {
			// Query returns a new read-only query builder for the {{.NameGo}} view.

			Query() query.Builder[model.{{.NameGo}}]
		}

		// {{.NameGo}}Repo returns the repository instance for the {{.NameGo}} view.
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

		// Query returns a new read-only query builder for the {{.NameGo}} view.
		func (r *{{.NameGoLower}}) Query() query.Builder[model.{{.NameGo}}] {
			return query.New{{.NameGo}}(r.db)
		}
	`

	data := map[string]any{
		"NameGo":      view.NameGo(),
		"NameGoLower": view.NameGoLower(),
	}

	return renderGoFileWithImports(
		b.fs.Writer(filepath.Join(def.PkgRepo, view.FileName())),
		def.PkgRepo, "repoView", tmpl, data,
		[]goImport{
			{Alias: "query", Path: b.relativePkgPath(def.PkgQuery)},
			{Alias: "model", Path: b.input.sourcePkgPath},
		},
	)
}
