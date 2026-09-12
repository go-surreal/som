package codegen

import (
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
)

// buildUnionRepoFile generates the read-only repository of a union. It queries
// all member tables at once and decodes every row into the member it belongs
// to. Writes are not exposed: they depend on the id type, hooks and features of
// a concrete member, so they stay with that member's repository.
func (b *build) buildUnionRepoFile(union *field.UnionTable) error {
	tmpl := `
		type {{.NameGo}}Repo interface {
			// Query returns a new query builder over all member tables of the {{.NameGo}} union.
			Query() query.Builder[model.{{.NameGo}}]
		}

		// {{.NameGo}}Repo returns the repository instance for the {{.NameGo}} union.
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

		// Query returns a new query builder over all member tables of the {{.NameGo}}
		// union: {{.Members}}. To create, update or delete a record, use the
		// repository of the member it belongs to.
		func (r *{{.NameGoLower}}) Query() query.Builder[model.{{.NameGo}}] {
			return query.New{{.NameGo}}(r.db)
		}
	`

	names := make([]string, len(union.Members))
	for i, member := range union.Members {
		names[i] = member.NameGo()
	}

	data := map[string]any{
		"NameGo":      union.NameGo(),
		"NameGoLower": union.NameGoLower(),
		"Members":     strings.Join(names, ", "),
	}

	file := newGoFile(def.PkgRepo,
		goImport{Alias: "query", Path: b.relativePkgPath(def.PkgQuery)},
		goImport{Alias: "model", Path: b.input.sourcePkgPath},
	)

	return file.render(
		b.fs.Writer(filepath.Join(def.PkgRepo, union.FileName())),
		"repoUnion", tmpl, data,
	)
}
