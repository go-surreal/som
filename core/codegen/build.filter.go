package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/util/fs"
)

type filterBuilder struct {
	*baseBuilder
}

func newFilterBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *filterBuilder {
	return &filterBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *filterBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildFile(edge); err != nil {
			return err
		}
	}

	for _, view := range b.views {
		if err := b.buildFile(view); err != nil {
			return err
		}
	}

	// Sinks are write-only, but their filter fields are generated so views
	// can project from them (the sink→view ingestion pattern). The refs are
	// only used to build view definitions, not to query the sink.
	for _, sink := range b.sinks {
		if err := b.buildFile(sink); err != nil {
			return err
		}
	}

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

// buildFile generates the filter accessors for a single table or object. Edges
// additionally get the traversal types for both of their ends, while all other
// tables get the type the traversal leads to.
//
// TODO: add record::exists filter function
// https://github.com/surrealdb/surrealdb/pull/4602
func (b *filterBuilder) buildFile(elem field.Element) error {
	tmpl := `
		var {{.NameGo}} = new{{.NameGo}}[model.{{.NameGo}}](lib.NewKey[model.{{.NameGo}}]())

		{{.NewFunc}}

		{{.StructType}}
		{{range .FieldFuncs}}
		{{.}}
		{{end}}
		{{- range .FieldExtras}}
		{{.}}
		{{end}}
		{{- if .Edge}}
		type {{.NameGoLower}}In[M any] struct {
			lib.Filter[M]
			key lib.Key[M]
		}

		func new{{.NameGo}}In[M any](key lib.Key[M]) {{.NameGoLower}}In[M] {
			return {{.NameGoLower}}In[M]{lib.KeyFilter(key), key}
		}

		func (i {{.NameGoLower}}In[M]) {{.Edge.OutNameGo}}(filters ...lib.Filter[model.{{.Edge.OutTableGo}}]) {{.Edge.OutTableLower}}Edges[M] {
			key := lib.EdgeIn(i.key, "{{.Edge.OutNameDB}}", filters)
			return {{.Edge.OutTableLower}}Edges[M]{lib.KeyFilter(key), key}
		}

		type {{.NameGoLower}}Out[M any] struct {
			lib.Filter[M]
			key lib.Key[M]
		}

		func new{{.NameGo}}Out[M any](key lib.Key[M]) {{.NameGoLower}}Out[M] {
			return {{.NameGoLower}}Out[M]{lib.KeyFilter(key), key}
		}

		func (o {{.NameGoLower}}Out[M]) {{.Edge.InNameGo}}(filters ...lib.Filter[model.{{.Edge.InTableGo}}]) {{.Edge.InTableLower}}Edges[M] {
			key := lib.EdgeOut(o.key, "{{.Edge.InNameDB}}", filters)
			return {{.Edge.InTableLower}}Edges[M]{lib.KeyFilter(key), key}
		}
		{{else}}
		// {{.NameGoLower}}Edges is the {{.NameGo}} as reached through a graph traversal.
		type {{.NameGoLower}}Edges[M any] struct {
			lib.Filter[M]
			lib.Key[M]
		}
		{{range .EdgeFuncs}}
		{{.}}
		{{end -}}
		{{end -}}
	`

	edge, isEdge := elem.(*field.EdgeTable)

	file := newGoFile(b.pkgName,
		goImport{Alias: "lib", Path: b.relativePkgPath(def.PkgLib)},
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)

	data := map[string]any{
		"NameGo":      elem.NameGo(),
		"NameGoLower": elem.NameGoLower(),
		"NewFunc":     file.decl(b.newFunc(elem)),
		"StructType":  file.decl(b.structType(elem)),
		"FieldFuncs":  b.fieldCodes(file, elem, (*field.CodeGen).FilterFunc),
		"FieldExtras": b.fieldCodes(file, elem, (*field.CodeGen).FilterExtra),
	}

	if isEdge {
		data["Edge"] = map[string]string{
			"InNameGo":      edge.In.NameGo(),
			"InNameDB":      edge.In.NameDatabase(),
			"InTableGo":     edge.In.Table().NameGo(),
			"InTableLower":  edge.In.Table().NameGoLower(),
			"OutNameGo":     edge.Out.NameGo(),
			"OutNameDB":     edge.Out.NameDatabase(),
			"OutTableGo":    edge.Out.Table().NameGo(),
			"OutTableLower": edge.Out.Table().NameGoLower(),
		}
	} else {
		data["EdgeFuncs"] = b.edgeFuncs(file, elem)
	}

	return file.render(
		b.fs.Writer(path.Join(b.path(), elem.FileName())),
		"filter", tmpl, data,
	)
}

// newFunc returns the constructor of the filter accessor struct.
func (b *filterBuilder) newFunc(elem field.Element) jen.Code {
	pkgLib := b.relativePkgPath(def.PkgLib)

	values := jen.Dict{
		jen.Id("Key"): jen.Id("key"),
	}

	for i, f := range elem.GetFields() {
		if code := f.CodeGen().FilterInit(b.fieldContext(elem, i)); code != nil {
			values[jen.Id(f.NameGo())] = code
		}
	}

	return jen.Func().Id("new" + elem.NameGo()).
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(def.TypeModel)).
		Id(elem.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id(elem.NameGoLower()).Types(def.TypeModel).Values(values)),
		)
}

// structType returns the type declaration of the filter accessor struct,
// holding one accessor per field of the given element.
func (b *filterBuilder) structType(elem field.Element) jen.Code {
	pkgLib := b.relativePkgPath(def.PkgLib)

	return jen.Type().Id(elem.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			// TODO: name clash with Key field! -> go1.23: type key_[M any] = lib.Key[M]
			g.Qual(pkgLib, "Key").Types(def.TypeModel)

			for i, f := range elem.GetFields() {
				if code := f.CodeGen().FilterDefine(b.fieldContext(elem, i)); code != nil {
					g.Add(code)
				}
			}
		})
}

// fieldCodes returns the given kind of filter declaration for all fields that
// need one.
func (b *filterBuilder) fieldCodes(
	file *goFile, elem field.Element,
	codeOf func(*field.CodeGen, field.Context) jen.Code,
) []string {
	var codes []string

	for i, f := range elem.GetFields() {
		if code := codeOf(f.CodeGen(), b.fieldContext(elem, i)); code != nil {
			codes = append(codes, file.decl(code))
		}
	}

	return codes
}

// edgeFuncs returns the traversal functions of all edge fields, bound to the
// type reached through a graph traversal.
func (b *filterBuilder) edgeFuncs(file *goFile, elem field.Element) []string {
	var funcs []string

	for i, f := range elem.GetFields() {
		if !isEdgeField(f) {
			continue
		}

		ctx := b.fieldContext(elem, i)
		ctx.Receiver = jen.Id(elem.NameGoLower() + "Edges").Types(def.TypeModel)

		if code := f.CodeGen().FilterFunc(ctx); code != nil {
			funcs = append(funcs, file.decl(code))
		}
	}

	return funcs
}

// isEdgeField reports whether the given field points to an edge, either
// directly or as a slice.
func isEdgeField(f field.Field) bool {
	if _, ok := f.(*field.Edge); ok {
		return true
	}

	if slice, ok := f.(*field.Slice); ok {
		_, ok := slice.Element().(*field.Edge)
		return ok
	}

	return false
}

func (b *filterBuilder) fieldContext(elem field.Element, index int) field.Context {
	return fieldContextFor(b.sourcePkgPath, b.basePkg, elem, index)
}
