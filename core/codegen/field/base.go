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

	Indexes() []parser.IndexInfo
	SearchInfo() *parser.SearchInfo
	NestedFields() []Field
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

// EdgeEnd is one end of an edge: either a single node table or a union of the
// node tables the end may point to.
type EdgeEnd interface {
	Field

	// TargetNameGo is the Go name of the element the end points to, which is
	// the node table or, for a multi-table end, the union of them.
	TargetNameGo() string
	TargetNameGoLower() string

	// TargetTables returns every node table the end may point to.
	TargetTables() []*NodeTable

	// RelationDatabase returns the endpoint type of the relation definition,
	// e.g. "square|triangle" for a multi-table end.
	RelationDatabase() string
}

// edgeEndHasTable reports whether the given end of an edge may point to the
// given table.
func edgeEndHasTable(end EdgeEnd, table Table) bool {
	for _, target := range end.TargetTables() {
		if tableEqual(table, target) {
			return true
		}
	}

	return false
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
	if dbName := f.source.DBName(); dbName != "" {
		return dbName
	}
	return f.ToDatabaseName(f.source.Name())
}

func (f *baseField) Indexes() []parser.IndexInfo {
	return f.source.Indexes()
}

func (f *baseField) SearchInfo() *parser.SearchInfo {
	return f.source.Search()
}

func (f *baseField) NestedFields() []Field {
	return nil
}
