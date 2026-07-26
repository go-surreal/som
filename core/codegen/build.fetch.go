package codegen

import (
	"path"

	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/util/fs"
)

type fetchBuilder struct {
	*baseBuilder
}

func newFetchBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *fetchBuilder {
	return &fetchBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *fetchBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

// fetchRelation describes a relation of a node that can be fetched.
type fetchRelation struct {
	NameGo      string
	NameDB      string
	TargetLower string

	// Kind names the relation within the doc comment.
	Kind string

	// SoftDelete marks relations pointing to a soft-deleted table, for which
	// the fetch behaviour is documented explicitly.
	SoftDelete bool
}

func (b *fetchBuilder) buildFile(node *field.NodeTable) error {
	tmpl := `
		var {{.NameGo}} = {{.NameGoLower}}[model.{{.NameGo}}]("")

		type {{.NameGoLower}}[M any] string

		func (n {{.NameGoLower}}[M]) fetch(M) {}
		{{range .Relations}}
		{{- if .SoftDelete}}
		// {{.NameGo}} returns a fetch accessor for the {{.NameDB}} {{.Kind}}.
		// Note: Soft-delete filtering does not apply to fetched relations.
		// All related records are returned regardless of their soft-delete status.
		{{- end}}
		func (n {{$.NameGoLower}}[M]) {{.NameGo}}() {{.TargetLower}}[M] {
			return {{.TargetLower}}[M](keyed(n, "{{.NameDB}}"))
		}
		{{end -}}
	`

	data := map[string]any{
		"NameGo":      node.NameGo(),
		"NameGoLower": node.NameGoLower(),
		"Relations":   fetchRelations(node),
	}

	file := newGoFile(b.pkgName,
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)

	return file.render(
		b.fs.Writer(path.Join(b.path(), node.FileName())),
		"fetch", tmpl, data,
	)
}

// fetchRelations returns the node and node slice fields of the given node,
// which are the fields a fetch can be applied to.
func fetchRelations(node *field.NodeTable) []fetchRelation {
	var relations []fetchRelation

	for _, fld := range node.GetFields() {
		var (
			target *field.NodeTable
			kind   string
		)

		switch fld := fld.(type) {
		case *field.Node:
			target, kind = fld.Table(), "relation"

		case *field.Slice:
			element, ok := fld.Element().(*field.Node)
			if !ok {
				continue
			}
			target, kind = element.Table(), "slice relation"

		default:
			continue
		}

		relations = append(relations, fetchRelation{
			NameGo:      fld.NameGo(),
			NameDB:      fld.NameDatabase(),
			TargetLower: target.NameGoLower(),
			Kind:        kind,
			SoftDelete:  target.Source != nil && target.Source.SoftDelete,
		})
	}

	return relations
}
