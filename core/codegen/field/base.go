package field

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

const (
	convTag     = "cbor"
	fnSuffixPtr = "Ptr"
)

// type Edge struct {
// 	Name   string
// 	In     Field
// 	Out    Field
// 	Fields []Field
// }

type ElemGetter func(name string) (Element, bool)

type Field interface {
	NameGo() string
	NameGoLower() string
	NameDatabase() string

	typeGo() jen.Code
	typeConv(ctx Context) jen.Code
	TypeDatabase() string

	SchemaStatements(table string, prefix string) []string

	CodeGen() *CodeGen
}

type Named interface {
	NameGo() string
	NameGoLower() string
	NameDatabase() string
}

type Element interface {
	Named

	FileName() string
	GetFields() []Field
	HasTimestamps() bool
}

type Table interface {
	Named

	FileName() string
	GetFields() []Field
}

type EnumModel string

func (m EnumModel) NameGo() string {
	return string(m)
}

func (m EnumModel) NameGoLower() string {
	return strcase.ToLowerCamel(string(m))
}

func (m EnumModel) NameDatabase() string {
	return strcase.ToSnake(string(m)) // TODO
}

func tableEqual(t1, t2 Table) bool {
	return t1.NameGo() == t2.NameGo()
}

type BuildConfig struct {
	SourcePkg      string
	TargetPkg      string
	ToDatabaseName func(base string) string
}

type baseField struct {
	*BuildConfig

	source parser.Field
}

func (f *baseField) ptr() jen.Code {
	if f.source.Pointer() {
		return jen.Op("*")
	}

	return jen.Empty()
}

func (f *baseField) omitEmptyIfPtr() string {
	if f.source.Pointer() {
		return ",omitempty"
	}

	return ""
}

// optionWrap wraps the given value in an option type if the field is a pointer.
func (f *baseField) optionWrap(val string) string {
	if f.source.Pointer() {
		return "option<" + val + ">"
	}

	return val
}

func (f *baseField) NameGo() string {
	return f.source.Name()
}

func (f *baseField) NameGoLower() string {
	return strings.ToLower(f.source.Name())
}

func (f *baseField) NameDatabase() string {
	return f.ToDatabaseName(f.source.Name())
}

// CollectPasswordPaths recursively collects database paths for all password fields.
// It traverses nested struct fields to find passwords at any depth.
func CollectPasswordPaths(fields []Field, prefix string) []string {
	var paths []string
	for _, f := range fields {
		switch field := f.(type) {
		case *Password:
			paths = append(paths, prefix+field.NameDatabase())
		case *Struct:
			nestedPrefix := prefix + field.NameDatabase() + "."
			paths = append(paths, CollectPasswordPaths(field.Table().GetFields(), nestedPrefix)...)
		}
	}
	return paths
}
