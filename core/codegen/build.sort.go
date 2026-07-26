package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/util/fs"
)

type sortBuilder struct {
	*baseBuilder
}

func newSortBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *sortBuilder {
	return &sortBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *sortBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
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

// buildFile generates the sort accessors for a single table or object. Only
// tables get an exported entry point, objects are reached through the field
// they are nested under.
func (b *sortBuilder) buildFile(elem field.Element) error {
	tmpl := `
		{{- if .IsTable}}
		var {{.NameGo}} = new{{.NameGo}}[model.{{.NameGo}}]("")
		{{end}}
		func new{{.NameGo}}[M any](key string) {{.NameGoLower}}[M] {
			return {{.InitLiteral}}
		}

		{{.StructType}}
		{{range .FieldFuncs}}
		{{.}}
		{{end -}}
	`

	_, isTable := elem.(*field.NodeTable)

	file := newGoFile(b.pkgName,
		goImport{Alias: "lib", Path: b.relativePkgPath(def.PkgLib)},
		goImport{Alias: "model", Path: b.sourcePkgPath},
	)

	data := map[string]any{
		"NameGo":      elem.NameGo(),
		"NameGoLower": elem.NameGoLower(),
		"IsTable":     isTable,
		"InitLiteral": file.code(b.initLiteral(elem)),
		"StructType":  file.decl(b.structType(elem)),
		"FieldFuncs":  b.fieldFuncs(file, elem),
	}

	return file.render(
		b.fs.Writer(path.Join(b.path(), elem.FileName())),
		"sort", tmpl, data,
	)
}

// initLiteral returns the composite literal initialising one sort accessor per
// field of the given element.
func (b *sortBuilder) initLiteral(elem field.Element) jen.Code {
	values := jen.Dict{
		jen.Id("key"): jen.Id("key"),
	}

	for i, f := range definedFields(elem) {
		if code := f.CodeGen().SortInit(b.fieldContext(elem, i)); code != nil {
			values[jen.Id(f.NameGo())] = code
		}
	}

	return jen.Id(elem.NameGoLower()).Types(def.TypeModel).Values(values)
}

// structType returns the type declaration of the accessor struct, holding one
// accessor per field of the given element.
func (b *sortBuilder) structType(elem field.Element) jen.Code {
	return jen.Type().Id(elem.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Id("key").String()

			for i, f := range definedFields(elem) {
				if code := f.CodeGen().SortDefine(b.fieldContext(elem, i)); code != nil {
					g.Add(code)
				}
			}
		})
}

// fieldFuncs returns the accessor functions of those fields that need one,
// like nested objects and slices.
func (b *sortBuilder) fieldFuncs(file *goFile, elem field.Element) []string {
	var funcs []string

	for _, f := range elem.GetFields() {
		if code := f.CodeGen().SortFunc(b.fieldContext(elem, 0)); code != nil {
			funcs = append(funcs, file.decl(code))
		}
	}

	return funcs
}

func (b *sortBuilder) fieldContext(elem field.Element, index int) field.Context {
	return fieldContextFor(b.sourcePkgPath, b.basePkg, elem, index)
}

// definedFields returns the fields an accessor struct is built from. In
// contrast to GetFields, these include the fields of the embedded som types.
func definedFields(elem field.Element) []field.Field {
	switch elem := elem.(type) {
	case *field.NodeTable:
		return elem.Fields

	case *field.DatabaseObject:
		return elem.Fields

	default:
		return elem.GetFields()
	}
}

func fieldContextFor(sourcePkg, targetPkg string, elem field.Element, index int) field.Context {
	ctx := field.Context{
		SourcePkg: sourcePkg,
		TargetPkg: targetPkg,
		Table:     elem,
	}

	// The fields of an array based complex ID are addressed by their position.
	if object, ok := elem.(*field.DatabaseObject); ok && object.IsArrayIndexed {
		ctx.ArrayIndex = &index
	}

	return ctx
}
