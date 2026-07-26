package codegen

import (
	"path/filepath"

	"github.com/go-surreal/som/core/codegen/def"
)

// buildWireFile generates the provider set for the configured wire package.
func (b *build) buildWireFile() error {
	tmpl := `
		var Providers = wire.NewSet(
			ProvideClient,
			wire.Bind(new(repo.Client), new(*repo.ClientImpl)),
			{{range .Repos}}
			Provide{{.NameGo}}Repo,
			{{- end}}
		)

		func ProvideClient(ctx context.Context, conf repo.Config) (*repo.ClientImpl, func(), error) {
			client, err := repo.NewClient(ctx, conf)
			if err != nil {
				return nil, nil, err
			}
			cleanup := func() {
				client.Close()
			}
			return client, cleanup, nil
		}
		{{range .Repos}}
		func Provide{{.NameGo}}Repo(client *repo.ClientImpl) repo.{{.NameGo}}Repo {
			return client.{{.NameGo}}Repo()
		}
		{{end -}}
	`

	var repos []repoRef

	for _, node := range b.input.nodes {
		repos = append(repos, repoRef{node.NameGo(), node.NameGoLower()})
	}

	data := map[string]any{
		"Repos": repos,
	}

	file := newGoFile(def.PkgSomWire,
		goImport{Path: "context"},
		goImport{Alias: "wire", Path: b.wirePackage},
		goImport{Alias: "repo", Path: b.relativePkgPath(def.PkgRepo)},
	)

	return file.render(
		b.fs.Writer(filepath.Join(def.PkgSomWire, "providers.go")),
		"wire", tmpl, data,
	)
}
