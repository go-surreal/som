package codegen

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/embed"
)

func renderSnippet(name, tmpl string, data map[string]any) (string, error) {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderGoFile(w io.Writer, pkgName, tmplName, tmplStr string, data map[string]any) error {
	snippet, err := renderSnippet(tmplName, tmplStr, data)
	if err != nil {
		return err
	}

	f := jen.NewFile(pkgName)
	f.PackageComment(string(embed.CodegenComment))
	f.Id(snippet)

	return f.Render(w)
}

type goImport struct {
	Alias string
	Path  string
}

func formatImportBlock(imports []goImport) string {
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	var lines []string

	for _, imp := range imports {
		if imp.Alias != "" {
			lines = append(lines, fmt.Sprintf("\t%s %q", imp.Alias, imp.Path))
		} else {
			lines = append(lines, fmt.Sprintf("\t%q", imp.Path))
		}
	}

	return "import (\n" + strings.Join(lines, "\n") + "\n)"
}
