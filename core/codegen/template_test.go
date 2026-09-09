package codegen

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestGoFileRender(t *testing.T) {
	file := newGoFile("repo",
		goImport{Path: "errors"},
		goImport{Alias: "model", Path: "example.com/app/model"},
		goImport{Alias: "unused", Path: "example.com/app/unused"},
	)

	data := map[string]any{
		"Check": file.code(jen.If(jen.Id("id").Op("==").Lit("")).Block(
			jen.Return(jen.Qual("errors", "New").Call(jen.Lit("empty id"))),
		)),
		"Helper": file.decl(jen.Func().Id("helper").Params().Qual("time", "Duration").Block(
			jen.Return(jen.Qual("time", "Second")),
		)),
	}

	tmpl := `
		func Create(id string) error {
			{{.Check}}
			_ = model.Thing{}
			return nil
		}

		{{.Helper}}
	`

	var out bytes.Buffer

	if err := file.render(&out, "test", tmpl, data); err != nil {
		t.Fatalf("render: %v", err)
	}

	got := out.String()

	for _, want := range []string{
		"package repo",
		`"errors"`,
		`model "example.com/app/model"`,
		`"time"`,
		"\tif id == \"\" {\n\t\treturn errors.New(\"empty id\")\n\t}",
		"func helper() time.Duration {",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rendered file does not contain %q:\n%s", want, got)
		}
	}

	// Candidates the code does not refer to must not be imported, and no
	// placeholder may survive the rendering.
	for _, unwanted := range []string{"unused", "__som_fragment"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("rendered file unexpectedly contains %q:\n%s", unwanted, got)
		}
	}
}

// TestGoFileRenderCollidingPackages covers two packages with the same name, as
// caused by the two supported UUID libraries. Jennifer aliases one of them, and
// both the import block and the rendered code have to agree on that alias.
func TestGoFileRenderCollidingPackages(t *testing.T) {
	file := newGoFile("field")

	data := map[string]any{
		"Google": file.code(jen.Var().Id("g").Qual("github.com/google/uuid", "UUID")),
		"Gofrs":  file.code(jen.Var().Id("f").Qual("github.com/gofrs/uuid", "UUID")),
	}

	var out bytes.Buffer

	if err := file.render(&out, "test", "func _() {\n{{.Google}}\n{{.Gofrs}}\n}", data); err != nil {
		t.Fatalf("render: %v", err)
	}

	got := out.String()

	if !strings.Contains(got, `"github.com/google/uuid"`) || !strings.Contains(got, `"github.com/gofrs/uuid"`) {
		t.Fatalf("both uuid packages should be imported:\n%s", got)
	}

	if !strings.Contains(got, "uuid1 ") {
		t.Fatalf("one of the uuid packages should be aliased:\n%s", got)
	}

	if strings.Count(got, "uuid1.UUID") != 1 {
		t.Errorf("the aliased package should be referenced by its alias:\n%s", got)
	}
}

// TestGoFileRenderPinnedAlias makes sure a package the template refers to by
// name keeps that name, even if a fragment pulls in a colliding one.
func TestGoFileRenderPinnedAlias(t *testing.T) {
	file := newGoFile("repo", goImport{Alias: "model", Path: "example.com/app/model"})

	data := map[string]any{
		"Other": file.code(jen.Var().Id("o").Qual("example.com/other/model", "Thing")),
	}

	var out bytes.Buffer

	if err := file.render(&out, "test", "func _() {\n_ = model.Thing{}\n{{.Other}}\n}", data); err != nil {
		t.Fatalf("render: %v", err)
	}

	got := out.String()

	if !strings.Contains(got, `model "example.com/app/model"`) {
		t.Errorf("the template package should keep its name:\n%s", got)
	}

	if !strings.Contains(got, `model1 "example.com/other/model"`) {
		t.Errorf("the colliding fragment package should be aliased:\n%s", got)
	}
}

func TestGoFileRenderInvalidCode(t *testing.T) {
	file := newGoFile("repo")

	err := file.render(&bytes.Buffer{}, "test", "func Broken( {", nil)
	if err == nil {
		t.Fatal("expected an error for code that does not parse")
	}
}

func TestRenderFragments(t *testing.T) {
	t.Run("statement is unwrapped and dedented", func(t *testing.T) {
		codes, _, err := renderFragments("repo", nil, []codeFragment{{
			code: jen.If(jen.Id("ok")).Block(jen.Return(jen.Nil())),
		}})
		if err != nil {
			t.Fatalf("renderFragments: %v", err)
		}

		want := "if ok {\n\treturn nil\n}"
		if codes[0] != want {
			t.Errorf("got %q, want %q", codes[0], want)
		}
	})

	t.Run("fragment with multiple declarations", func(t *testing.T) {
		codes, _, err := renderFragments("repo", nil, []codeFragment{{
			isDecl: true,
			code: jen.Type().Id("a").String().Line().Line().
				Type().Id("b").String(),
		}})
		if err != nil {
			t.Fatalf("renderFragments: %v", err)
		}

		for _, want := range []string{"type a string", "type b string"} {
			if !strings.Contains(codes[0], want) {
				t.Errorf("got %q, want it to contain %q", codes[0], want)
			}
		}
	})

	t.Run("declaration keeps its doc comment", func(t *testing.T) {
		codes, _, err := renderFragments("repo", nil, []codeFragment{{
			isDecl: true,
			code:   jen.Comment("// Thing is a thing.").Line().Type().Id("Thing").String(),
		}})
		if err != nil {
			t.Fatalf("renderFragments: %v", err)
		}

		if !strings.HasPrefix(codes[0], "// Thing is a thing.") {
			t.Errorf("got %q, want it to start with the doc comment", codes[0])
		}
	})

	t.Run("no fragments", func(t *testing.T) {
		codes, imports, err := renderFragments("repo", nil, nil)
		if err != nil {
			t.Fatalf("renderFragments: %v", err)
		}

		if codes != nil || imports != nil {
			t.Errorf("got %v / %v, want nothing", codes, imports)
		}
	})
}

func TestMergeImports(t *testing.T) {
	t.Run("duplicates are removed", func(t *testing.T) {
		merged, err := mergeImports([]goImport{
			{Alias: "errors", Path: "errors"},
			{Path: "errors"},
		})
		if err != nil {
			t.Fatalf("mergeImports: %v", err)
		}

		if len(merged) != 1 {
			t.Fatalf("got %d imports, want 1: %v", len(merged), merged)
		}

		// The last one wins, as that is the alias the template relies on.
		if merged[0].Alias != "" {
			t.Errorf("got alias %q, want none", merged[0].Alias)
		}
	})

	t.Run("name clash is rejected", func(t *testing.T) {
		_, err := mergeImports([]goImport{
			{Alias: "uuid", Path: "github.com/google/uuid"},
			{Alias: "uuid", Path: "github.com/gofrs/uuid"},
		})
		if err == nil {
			t.Fatal("expected an error for two packages sharing a name")
		}
	})
}

func TestUsedImports(t *testing.T) {
	candidates := []goImport{
		{Alias: "conv", Path: "example.com/app/conv"},
		{Alias: "models", Path: "example.com/models"},
		{Alias: "query", Path: "example.com/app/query"},
		{Path: "errors"},
	}

	body := `
		// NewThing creates a new query builder for Thing models.
		func NewThing(fetch ...query.Fetch) error {
			_ = conv.Thing{}
			_ = r.query.Table
			return errors.New("no models. here")
		}
	`

	var names []string
	for _, imp := range usedImports(body, candidates) {
		names = append(names, imp.name())
	}

	// query is only referenced through a variadic parameter, models only in
	// prose, and r.query is a selector rather than a package.
	want := []string{"conv", "query", "errors"}
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Errorf("got %v, want %v", names, want)
	}
}

func TestStripCommentsAndStrings(t *testing.T) {
	got := stripCommentsAndStrings(`a := "text" // comment` + "\nb := c")

	if strings.Contains(got, "text") || strings.Contains(got, "comment") {
		t.Errorf("got %q, want strings and comments removed", got)
	}

	if !strings.Contains(got, "b := c") {
		t.Errorf("got %q, want the code to be kept", got)
	}
}

func TestFormatImportBlock(t *testing.T) {
	t.Run("single import", func(t *testing.T) {
		got := formatImportBlock([]goImport{{Alias: "model", Path: "example.com/model"}})

		want := `import model "example.com/model"`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("sorted by path", func(t *testing.T) {
		got := formatImportBlock([]goImport{
			{Alias: "model", Path: "example.com/model"},
			{Path: "errors"},
		})

		want := "import (\n\t\"errors\"\n\tmodel \"example.com/model\"\n)"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestFormatGoComment(t *testing.T) {
	got := formatGoComment("First line.\n\nSecond line.")

	want := "// First line.\n//\n// Second line."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestJoinStatements(t *testing.T) {
	codes, _, err := renderFragments("repo", nil, []codeFragment{{
		code: joinStatements([]jen.Code{
			jen.Id("a").Op(":=").Lit(1),
			jen.Id("b").Op(":=").Lit(2),
		}),
	}})
	if err != nil {
		t.Fatalf("renderFragments: %v", err)
	}

	want := "a := 1\nb := 2"
	if codes[0] != want {
		t.Errorf("got %q, want %q", codes[0], want)
	}
}
