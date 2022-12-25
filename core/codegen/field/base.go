package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/parser"
)

const (
	funcBuildDatabaseID = "buildDatabaseID"
	funcParseDatabaseID = "parseDatabaseID"
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
	typeConv() jen.Code
	// typeDatabase() string

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
}

type Table interface {
	Named

	FileName() string
	GetFields() []Field
}

type Model string // TODO: just use table instead?

func (m Model) NameGo() string {
	return string(m)
}

func (m Model) NameGoLower() string {
	return strcase.ToLowerCamel(string(m))
}

func (m Model) NameDatabase() string {
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

func (f *baseField) NameGo() string {
	return f.source.Name()
}

func (f *baseField) NameGoLower() string {
	return strcase.ToLowerCamel(f.NameGo())
}

func (f *baseField) NameDatabase() string {
	return f.ToDatabaseName(f.source.Name())
}
