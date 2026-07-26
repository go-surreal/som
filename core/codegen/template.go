package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/embed"
)

// goFile assembles a generated Go file from a template.
//
// Templates hold the structure of the generated code, while jennifer builds
// those fragments that depend on the type of a field. Fragments are rendered
// through code(), which also records the packages they reference, so that the
// import block can be derived from them the same way jennifer would.
type goFile struct {
	pkg string

	// imports lists the packages the template itself may reference. Only those
	// actually referenced by the rendered code are added to the file.
	imports []goImport

	// fragments holds the jennifer fragments interpolated into the template.
	fragments []codeFragment
}

// codeFragment is a piece of jennifer built code interpolated into a template.
type codeFragment struct {
	code jen.Code

	// isDecl marks declarations - types and functions - which, in contrast to
	// statements, are only valid at the top level of a file.
	isDecl bool
}

func newGoFile(pkg string, imports ...goImport) *goFile {
	return &goFile{
		pkg:     pkg,
		imports: imports,
	}
}

// code records a jennifer statement or expression for interpolation into the
// template and returns a placeholder for it. The fragments are rendered once
// the template has been executed, as their package aliases can only be
// resolved when all of them are known.
func (f *goFile) code(fragment jen.Code) string {
	return f.add(codeFragment{code: fragment})
}

// decl records a jennifer declaration - a type or function - for interpolation
// into the template and returns a placeholder for it.
func (f *goFile) decl(fragment jen.Code) string {
	return f.add(codeFragment{code: fragment, isDecl: true})
}

func (f *goFile) add(fragment codeFragment) string {
	f.fragments = append(f.fragments, fragment)

	return fmt.Sprintf(fragmentPlaceholder, len(f.fragments)-1)
}

// render executes the given template and writes the resulting file.
func (f *goFile) render(w io.Writer, name, tmpl string, data map[string]any) error {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", name, err)
	}

	var body bytes.Buffer

	// Note: this is what records the fragments of the template functions,
	// so they can only be rendered afterwards.
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("could not execute template %s: %w", name, err)
	}

	content, imports, err := f.resolveFragments(body.String())
	if err != nil {
		return fmt.Errorf("could not render fragments of %s: %w", name, err)
	}

	if len(imports) > 0 {
		content = formatImportBlock(imports) + "\n\n" + content
	}

	out := jen.NewFile(f.pkg)
	out.PackageComment(string(embed.CodegenComment))
	out.Id(content)

	return out.Render(w)
}

// resolveFragments replaces the fragment placeholders of the rendered template
// by the rendered fragments and returns the imports of the resulting file.
func (f *goFile) resolveFragments(body string) (string, []goImport, error) {
	fragments, imports, err := renderFragments(f.pkg, f.imports, f.fragments)
	if err != nil {
		return "", nil, err
	}

	for i, fragment := range fragments {
		body = strings.ReplaceAll(body, fmt.Sprintf(fragmentPlaceholder, i), fragment)
	}

	// The template packages are only imported if the code refers to them, while
	// the fragment packages are always needed by definition. They are merged
	// last, so that the alias the template relies on wins.
	merged, err := mergeImports(append(imports, usedImports(body, f.imports)...))
	if err != nil {
		return "", nil, err
	}

	return body, merged, nil
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

// fragmentPlaceholder marks the position of a fragment within a rendered
// template, until the fragment itself is rendered.
const (
	fragmentPlaceholder = "__som_fragment_%d__"

	// fragmentFuncName is the throwaway function a statement fragment is
	// wrapped into, to make it valid at the top level of a file.
	fragmentFuncName = "fragment"
)

// renderFragments renders the fragments of a file and returns their code
// together with the imports they need.
//
// The aliases of the referenced packages are resolved across all fragments
// first, so that a fragment cannot end up referring to a package by a name that
// another fragment - or the template itself - already uses.
func renderFragments(pkg string, pinned []goImport, fragments []codeFragment) ([]string, []goImport, error) {
	if len(fragments) == 0 {
		return nil, nil, nil
	}

	imports, err := fragmentImports(pkg, pinned, fragments)
	if err != nil {
		return nil, nil, err
	}

	codes := make([]string, 0, len(fragments))

	for _, fragment := range fragments {
		code, err := renderFragment(pkg, imports, fragment)
		if err != nil {
			return nil, nil, err
		}

		codes = append(codes, code)
	}

	return codes, imports, nil
}

// fragmentImports returns the packages jennifer resolves for all fragments of a
// file. This is what makes the imports of a fragment known without the builders
// having to know upfront which packages a field type may pull in.
func fragmentImports(pkg string, pinned []goImport, fragments []codeFragment) ([]goImport, error) {
	file := newFragmentFile(pkg, pinned)

	for _, fragment := range fragments {
		addFragment(file, fragment)
	}

	source, err := renderFile(file)
	if err != nil {
		return nil, err
	}

	parsed, err := parser.ParseFile(token.NewFileSet(), "", source, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("could not parse code fragments: %w", err)
	}

	var imports []goImport

	for _, imp := range parsed.Imports {
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		}

		imports = append(imports, goImport{
			Alias: alias,
			Path:  strings.Trim(imp.Path.Value, `"`),
		})
	}

	return mergeImports(imports)
}

// renderFragment renders a single fragment by rendering it as a file, of which
// everything but the declarations is cut off again. Statements are wrapped into
// a throwaway function to make them valid at the top level of that file.
func renderFragment(pkg string, pinned []goImport, fragment codeFragment) (string, error) {
	file := newFragmentFile(pkg, pinned)
	addFragment(file, fragment)

	source, err := renderFile(file)
	if err != nil {
		return "", err
	}

	parsed, err := parser.ParseFile(token.NewFileSet(), "", source, parser.SkipObjectResolution)
	if err != nil {
		return "", fmt.Errorf("could not parse code fragment: %w", err)
	}

	for _, decl := range parsed.Decls {
		if generic, ok := decl.(*ast.GenDecl); ok && generic.Tok == token.IMPORT {
			continue
		}

		return declCode(source, decl), nil
	}

	return "", nil
}

func newFragmentFile(pkg string, pinned []goImport) *jen.File {
	file := jen.NewFile(pkg)

	// The template refers to these packages by a fixed name, so they must keep
	// it. Colliding packages are aliased by jennifer instead.
	for _, imp := range pinned {
		file.ImportAlias(imp.Path, imp.name())
	}

	return file
}

func addFragment(file *jen.File, fragment codeFragment) {
	if fragment.isDecl {
		file.Add(fragment.code)
		return
	}

	file.Func().Id(fragmentFuncName).Params().Block(fragment.code)
}

func renderFile(file *jen.File) ([]byte, error) {
	var rendered bytes.Buffer

	if err := file.Render(&rendered); err != nil {
		return nil, fmt.Errorf("could not render code fragment: %w", err)
	}

	return rendered.Bytes(), nil
}

// declCode cuts the code of the given declaration - and everything following
// it, as a fragment may hold more than one - out of the rendered file. The body
// of the throwaway function is unwrapped, as it holds a fragment that is not a
// declaration itself.
func declCode(source []byte, decl ast.Decl) string {
	if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == fragmentFuncName {
		body := string(source[fn.Body.Lbrace : fn.Body.Rbrace-1])

		return dedent(strings.Trim(body, "\n"))
	}

	start := decl.Pos()

	switch decl := decl.(type) {
	case *ast.GenDecl:
		if decl.Doc != nil {
			start = decl.Doc.Pos()
		}
	case *ast.FuncDecl:
		if decl.Doc != nil {
			start = decl.Doc.Pos()
		}
	}

	return strings.TrimRight(string(source[start-1:]), "\n")
}

// dedent removes one level of indentation, which a fragment gained by being
// rendered as the body of a function.
func dedent(code string) string {
	lines := strings.Split(code, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, "\t")
	}

	return strings.Join(lines, "\n")
}

// mergeImports removes duplicate imports and rejects name clashes.
func mergeImports(imports []goImport) ([]goImport, error) {
	byName := make(map[string]goImport, len(imports))

	for _, imp := range imports {
		if existing, ok := byName[imp.name()]; ok && existing.Path != imp.Path {
			return nil, fmt.Errorf("packages %q and %q share the name %q", existing.Path, imp.Path, imp.name())
		}

		byName[imp.name()] = imp
	}

	merged := make([]goImport, 0, len(byName))
	for _, imp := range byName {
		merged = append(merged, imp)
	}

	return merged, nil
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

// renderStatement executes a template for generated content that is not Go
// code, like the statements of the database schema.
func renderStatement(tmpl *template.Template, data any) (string, error) {
	var rendered strings.Builder

	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("could not execute template %s: %w", tmpl.Name(), err)
	}

	return rendered.String(), nil
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
