package field

import (
	"fmt"
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

// SchemaStatements returns the DEFINE FIELD statements for this field.
// This default implementation works for simple fields. Complex fields
// like Struct and Slice override this to handle nested fields.
func (f *baseField) SchemaStatements(table, prefix string, fieldType string) []string {
	if fieldType == "" {
		return nil
	}

	statement := fmt.Sprintf(
		"DEFINE FIELD %s ON TABLE %s TYPE %s;",
		prefix+f.NameDatabase(), table, fieldType,
	)

	return []string{statement}
}

// schemaStatement is a helper that generates a single DEFINE FIELD statement
// with an optional extension (ASSERT, VALUE, etc.).
func (f *baseField) schemaStatement(table, prefix, fieldType, extend string) string {
	if extend != "" {
		fieldType = fieldType + " " + extend
	}
	return fmt.Sprintf(
		"DEFINE FIELD %s ON TABLE %s TYPE %s;",
		prefix+f.NameDatabase(), table, fieldType,
	)
}
