package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/util/fs"
)

// accessorBuilder generates the key based accessor structs of the sort and
// field packages. Both hold one accessor per field of a table or object and
// differ only in the kind of accessor they are built from.
type accessorBuilder struct {
	*baseBuilder

	// name identifies the generated accessors in error messages.
	name string

	// initOf, defineOf and funcOf return the code of a single field: its
	// initialisation, its struct field and its accessor function.
	initOf   accessorCode
	defineOf accessorCode
	funcOf   accessorCode
}

type accessorCode func(*field.CodeGen, field.Context) jen.Code

func newSortBuilder(input *input, fs *fs.FS, basePkg, pkgName string) builder {
	return &accessorBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
		name:        "sort",
		initOf:      (*field.CodeGen).SortInit,
		defineOf:    (*field.CodeGen).SortDefine,
		funcOf:      (*field.CodeGen).SortFunc,
	}
}

func newFieldBuilder(input *input, fs *fs.FS, basePkg, pkgName string) builder {
	return &accessorBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
		name:        "field",
		initOf:      (*field.CodeGen).FieldInit,
		defineOf:    (*field.CodeGen).FieldDefine,
		funcOf:      (*field.CodeGen).FieldFunc,
	}
}

func (b *accessorBuilder) build() error {
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

// buildFile generates the accessors for a single table or object. Only tables
// get an exported entry point, objects are reached through the field they are
// nested under.
func (b *accessorBuilder) buildFile(elem field.Element) error {
	tmpl := `
		{{- if .IsTable}}
		var {{.NameGo}} = new{{.NameGo}}[model.{{.NameGo}}]("")
		{{end}}
		func new{{.NameGo}}[M any](key string) {{.NameGoLower}}[M] {
			return {{.InitLiteral}}
		}

		{{.StructType}}
		{{range $fn := .FieldFuncs}}
		{{$fn}}
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
		b.name, tmpl, data,
	)
}

// initLiteral returns the composite literal initialising one accessor per
// field of the given element.
func (b *accessorBuilder) initLiteral(elem field.Element) jen.Code {
	values := jen.Dict{
		jen.Id("key"): jen.Id("key"),
	}

	for i, f := range definedFields(elem) {
		if code := b.initOf(f.CodeGen(), b.fieldContext(elem, i)); code != nil {
			values[jen.Id(f.NameGo())] = code
		}
	}

	return jen.Id(elem.NameGoLower()).Types(def.TypeModel).Values(values)
}

// structType returns the type declaration of the accessor struct, holding one
// accessor per field of the given element.
func (b *accessorBuilder) structType(elem field.Element) jen.Code {
	return jen.Type().Id(elem.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Id("key").String()

			for i, f := range definedFields(elem) {
				if code := b.defineOf(f.CodeGen(), b.fieldContext(elem, i)); code != nil {
					g.Add(code)
				}
			}
		})
}

// fieldFuncs returns the accessor functions of those fields that need one,
// like nested objects and slices.
func (b *accessorBuilder) fieldFuncs(file *goFile, elem field.Element) []string {
	var funcs []string

	for _, f := range elem.GetFields() {
		if code := b.funcOf(f.CodeGen(), b.fieldContext(elem, 0)); code != nil {
			funcs = append(funcs, file.decl(code))
		}
	}

	return funcs
}

func (b *accessorBuilder) fieldContext(elem field.Element, index int) field.Context {
	return fieldContextFor(b.sourcePkgPath, b.basePkg, elem, index)
}
