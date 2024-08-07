package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
	"strings"
)

var (
	typeModel = jen.Id("T") // TODO: rename to M
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

// optionWrap wraps the given value in an option type if the field is a pointer.
func (f *baseField) optionWrap(val string) string {
	if f.source.Pointer() {
		return "option<" + val + " | null>"
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
