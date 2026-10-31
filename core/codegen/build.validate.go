package codegen

import (
	"path/filepath"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/util/fs"
)

// validation knows which model types take part in validation and emits the
// walkers checking them.
//
// A type takes part if it provides a "Validate() error" method itself, or if
// any of its fields holds a type that does. Nothing else is touched, so a
// model without validation keeps the exact code it had before.
//
// Validation is a write-time concern only: a record stored before a rule
// existed would otherwise be unreadable, and hence unfixable.
type validation struct {
	*baseBuilder

	// structs holds the nested struct types by their Go name.
	structs map[string]*field.DatabaseObject

	// needs caches, per struct type, whether it needs a walker at all.
	needs map[string]bool

	// pending guards against a struct type that (indirectly) contains itself.
	pending map[string]bool
}

func (b *build) newValidation() *validation {
	return newValidation(b.input, b.fs, b.basePkg(), def.PkgValidate)
}

func newValidation(input *input, fs *fs.FS, basePkg, pkgName string) *validation {
	v := &validation{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
		structs:     make(map[string]*field.DatabaseObject, len(input.objects)),
		needs:       make(map[string]bool),
		pending:     make(map[string]bool),
	}

	for _, object := range input.objects {
		v.structs[object.Name] = object
	}

	return v
}

// Node reports whether writes of the given node are validated.
func (v *validation) Node(node *field.NodeTable) bool {
	return node.HasValidate || v.anyField(node.Fields)
}

// Edge reports whether writes of the given edge are validated.
func (v *validation) Edge(edge *field.EdgeTable) bool {
	return edge.HasValidate || v.anyField(edge.Fields)
}

// Sink reports whether writes of the given sink are validated.
func (v *validation) Sink(sink *field.SinkTable) bool {
	return sink.HasValidate || v.anyField(sink.Fields)
}

// object reports whether the nested struct type with the given name needs a
// walker of its own.
func (v *validation) object(name string) bool {
	if needs, ok := v.needs[name]; ok {
		return needs
	}

	object, ok := v.structs[name]
	if !ok {
		return false
	}

	// A struct type that contains itself is reported as not needing validation
	// while it is being examined, so that the walk terminates.
	if v.pending[name] {
		return false
	}
	v.pending[name] = true

	needs := object.HasValidate || v.anyField(object.Fields)

	delete(v.pending, name)
	v.needs[name] = needs

	return needs
}

func (v *validation) anyField(fields []field.Field) bool {
	for _, f := range fields {
		if v.field(f) {
			return true
		}
	}

	return false
}

// field reports whether the given field contributes anything to validate.
//
// Links to other records are never followed: they are records of their own,
// written through their own repository, and a link that was not fetched holds
// nothing but its id.
func (v *validation) field(f field.Field) bool {
	if f.SkipValidate() {
		return false
	}

	switch f := f.(type) {
	case *field.Node, *field.Edge, *field.Union, *field.ID, *field.ComplexID:
		return false

	case *field.Struct:
		return v.object(f.Table().NameGo())

	case *field.Slice:
		return v.field(f.Element())
	}

	return f.HasValidate()
}

func (v *validation) build() error {
	var decls []string

	file := newGoFile(v.pkgName,
		goImport{Alias: "internal", Path: v.relativePkgPath(def.PkgInternal)},
		goImport{Alias: "model", Path: v.sourcePkgPath},
	)

	for _, node := range v.nodes {
		if v.Node(node) {
			decls = append(decls, v.tableWalker(file, node.NameGo(), node.Fields, node.HasValidate)...)
		}
	}

	for _, edge := range v.edges {
		if v.Edge(edge) {
			decls = append(decls, v.tableWalker(file, edge.NameGo(), edge.Fields, edge.HasValidate)...)
		}
	}

	for _, sink := range v.sinks {
		if v.Sink(sink) {
			decls = append(decls, v.tableWalker(file, sink.NameGo(), sink.Fields, sink.HasValidate)...)
		}
	}

	for _, object := range v.objects {
		if v.object(object.Name) {
			decls = append(decls, v.walkerFunc(file, object.Name, object.Fields, object.HasValidate))
		}
	}

	// Without a single model to check, the package is not generated at all.
	if len(decls) == 0 {
		return nil
	}

	return file.render(
		v.fs.Writer(filepath.Join(v.path(), "som.validate.go")),
		"validate", "{{range .Decls}}{{.}}\n\n{{end}}",
		map[string]any{"Decls": decls},
	)
}

// tableWalker emits the entry point of a table and the function doing the
// actual walk, which nested types call as well.
func (v *validation) tableWalker(file *goFile, name string, fields []field.Field, own bool) []string {
	entry := jen.Comment(name + " checks the " + name + " model and every value below it that").Line().
		Comment("validates itself. It is called by the repository before a write.").Line().
		Func().Id(name).Params(jen.Id("node").Op("*").Qual(v.sourcePkgPath, name)).Error().Block(
		jen.Return(jen.Id(walkerName(name)).Call(jen.Id("node"), jen.Lit(""))),
	)

	return []string{file.decl(entry), v.walkerFunc(file, name, fields, own)}
}

// walkerFunc emits the function validating one type: first every field that
// holds something to validate, then the type itself. Leaves are checked before
// their root, so that the most specific error wins.
func (v *validation) walkerFunc(file *goFile, name string, fields []field.Field, own bool) string {
	var stmts []jen.Code

	for _, f := range fields {
		if !v.field(f) {
			continue
		}

		value := jen.Id("val").Dot(f.NameGo())
		path := jen.Qual(v.internalPkg(), "FieldPath").Call(jen.Id("path"), jen.Lit(f.NameGo()))

		stmts = append(stmts, v.checks(f, value, path, 0)...)
	}

	if own {
		stmts = append(stmts, v.validateCall(jen.Id("val"), jen.Id("path")))
	}

	stmts = append(stmts, jen.Return(jen.Nil()))

	fn := jen.Comment(walkerName(name)+" validates a "+name+" value at the given model path.").Line().
		Func().Id(walkerName(name)).
		Params(
			jen.Id("val").Op("*").Qual(v.sourcePkgPath, name),
			jen.Id("path").String(),
		).
		Error().
		Block(stmts...)

	return file.decl(fn)
}

// checks emits the statements validating a single value, which for a slice is
// every one of its elements and for a nested struct the walker of that struct.
// The depth is the number of slices already entered, used to keep the index
// variables apart.
func (v *validation) checks(f field.Field, value, path *jen.Statement, depth int) []jen.Code {
	switch f := f.(type) {

	case *field.Slice:
		index := jen.Id("i" + strconv.Itoa(depth))

		// A pointer to a slice has to be dereferenced before it can be ranged
		// over, as well as to address one of its elements.
		slice := value.Clone()
		if f.IsPointer() {
			slice = jen.Parens(jen.Op("*").Add(value.Clone()))
		}

		elemValue := slice.Clone().Index(index.Clone())
		elemPath := jen.Qual(v.internalPkg(), "IndexPath").Call(path.Clone(), index.Clone())

		loop := jen.For(index.Clone().Op(":=").Range().Add(slice.Clone())).
			Block(v.checks(f.Element(), elemValue, elemPath, depth+1)...)

		return []jen.Code{v.nilGuard(f, value, loop)}

	case *field.Struct:
		target := jen.Op("&").Add(value.Clone())
		if f.IsPointer() {
			target = value.Clone()
		}

		call := jen.If(
			jen.Err().Op(":=").Id(walkerName(f.Table().NameGo())).Call(target, path.Clone()),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Return(jen.Err()),
		)

		return []jen.Code{v.nilGuard(f, value, call)}
	}

	return []jen.Code{v.nilGuard(f, value, v.validateCall(value.Clone(), path.Clone()))}
}

// validateCall emits the call to the Validate method of a value, wrapping a
// rejection into a ValidationError carrying the path of that value.
func (v *validation) validateCall(value, path *jen.Statement) jen.Code {
	return jen.If(
		jen.Err().Op(":=").Add(value.Clone()).Dot("Validate").Call(),
		jen.Err().Op("!=").Nil(),
	).Block(
		jen.Return(jen.Qual(v.internalPkg(), "NewValidationError").Call(path.Clone(), jen.Err())),
	)
}

// nilGuard wraps the given statement into a nil check if the field is a
// pointer, since an absent value has nothing to validate.
func (v *validation) nilGuard(f field.Field, value *jen.Statement, stmt jen.Code) jen.Code {
	if !f.IsPointer() {
		return stmt
	}

	return jen.If(value.Clone().Op("!=").Nil()).Block(stmt)
}

func (v *validation) internalPkg() string {
	return v.relativePkgPath(def.PkgInternal)
}

// walkerName returns the name of the function validating the given type.
func walkerName(name string) string {
	return "validate" + name
}
