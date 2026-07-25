package codegen

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"regexp"
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

// renderGoFileWithImports renders the given template and prepends an import
// block containing those candidates that the rendered code actually references.
// This mirrors the automatic import handling of jennifer for template based
// generation, where the used packages are not known upfront.
func renderGoFileWithImports(
	w io.Writer, pkgName, tmplName, tmplStr string,
	data map[string]any, candidates []goImport,
) error {
	body, err := renderSnippet(tmplName, tmplStr, data)
	if err != nil {
		return err
	}

	imports := usedImports(body, candidates)

	f := jen.NewFile(pkgName)
	f.PackageComment(string(embed.CodegenComment))
	f.Id(formatImportBlock(imports) + "\n\n" + body)

	return f.Render(w)
}

type goImport struct {
	Alias string
	Path  string
}

func (i goImport) name() string {
	if i.Alias != "" {
		return i.Alias
	}
	return path.Base(i.Path)
}

// usedImports filters the given candidates down to those referenced by the
// rendered code via their package name (e.g. "conv." for the conv package).
func usedImports(body string, candidates []goImport) []goImport {
	code := stripCommentsAndStrings(body)

	var used []goImport

	for _, imp := range candidates {
		pattern := regexp.MustCompile(`(^|[^\w.])` + regexp.QuoteMeta(imp.name()) + `\.`)

		if pattern.MatchString(code) {
			used = append(used, imp)
		}
	}

	return used
}

// commentsAndStrings matches line comments as well as string literals, as
// neither of them can contain an actual package reference.
var commentsAndStrings = regexp.MustCompile(`//[^\n]*|"(\\.|[^"\\])*"`)

func stripCommentsAndStrings(body string) string {
	return commentsAndStrings.ReplaceAllString(body, "")
}

func formatImportBlock(imports []goImport) string {
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	var lines []string

	for _, imp := range imports {
		if imp.Alias != "" {
			lines = append(lines, fmt.Sprintf("%s %q", imp.Alias, imp.Path))
		} else {
			lines = append(lines, fmt.Sprintf("%q", imp.Path))
		}
	}

	if len(lines) == 1 {
		return "import " + lines[0]
	}

	return "import (\n\t" + strings.Join(lines, "\n\t") + "\n)"
}

func formatGoComment(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")

	out := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			out = append(out, "//")
		} else {
			out = append(out, "// "+line)
		}
	}

	return strings.Join(out, "\n")
}

// renderCode renders a jennifer code fragment into its string representation,
// so that it can be embedded into a template. Package qualifiers are rendered
// using the last element of their path, matching the aliases used within the
// generated import blocks.
func renderCode(code jen.Code) string {
	return fmt.Sprintf("%#v", code)
}
