package codegen

import (
	"path/filepath"

	"github.com/go-surreal/som/core/codegen/def"
)

const filenameInterfaces = "som.interfaces.go"

// repoRef describes one of the repositories exposed by the client.
type repoRef struct {
	NameGo      string
	NameGoLower string
}

// buildInterfaceFile generates the client interface and its implementation,
// holding one cached repository instance per table.
func (b *build) buildInterfaceFile() error {
	tmpl := `
		type Client interface {
			{{- range .Repos}}
			{{.NameGo}}Repo() {{.NameGo}}Repo
			{{- end}}
			Raw(ctx context.Context, query string, params som.Params) (*som.RawResult, error)
			ApplySchema(ctx context.Context) error
			Close()
		}

		type ClientImpl struct {
			db *dbConn
			mu sync.Mutex
			{{- range .Repos}}
			{{.NameGoLower}}Repo *{{.NameGoLower}}
			{{- end}}
		}

		// expiryTables lists tables with a configured expiry, purged in the background.
		var expiryTables = []string{ {{- range $i, $t := .ExpiryTables}}{{if $i}}, {{end}}"{{$t}}"{{end -}} }
	`

	var repos []repoRef
	var expiryTables []string

	for _, node := range b.input.nodes {
		repos = append(repos, repoRef{node.NameGo(), node.NameGoLower()})

		if node.Source.Expiry {
			expiryTables = append(expiryTables, node.NameDatabase())
		}
	}

	for _, view := range b.input.views {
		repos = append(repos, repoRef{view.NameGo(), view.NameGoLower()})
	}

	for _, sink := range b.input.sinks {
		repos = append(repos, repoRef{sink.NameGo(), sink.NameGoLower()})
	}

	data := map[string]any{
		"Repos":        repos,
		"ExpiryTables": expiryTables,
	}

	file := newGoFile(def.PkgRepo,
		goImport{Path: "context"},
		goImport{Path: "sync"},
		goImport{Alias: "som", Path: b.relativePkgPath()},
	)

	return file.render(
		b.fs.Writer(filepath.Join(def.PkgRepo, filenameInterfaces)),
		"repoClient", tmpl, data,
	)
}
