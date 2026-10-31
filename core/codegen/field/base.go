package field

import (
	"fmt"
	"strings"
	"text/template"

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

	// HasValidate reports whether the type of the field provides a
	// "Validate() error" method to be called before a write.
	HasValidate() bool

	// SkipValidate reports whether the field opted out of validation.
	SkipValidate() bool

	// IsPointer reports whether the model field is a pointer.
	IsPointer() bool
}

// elementTyped is implemented by fields whose database type differs when they
// occur as a slice element instead of as a field of their own.
type elementTyped interface {
	ElementTypeDatabase() string
}

// elementTypeDatabase returns the database type of the given field as it occurs
// within a slice.
func elementTypeDatabase(f Field) string {
	if elem, ok := f.(elementTyped); ok {
		return elem.ElementTypeDatabase()
	}

	return f.TypeDatabase()
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

	// Target is the element the end points to: the node table or, for a
	// multi-table end, the union of them.
	Target() Named

	// TargetTables returns every node table the end may point to.
	TargetTables() []*NodeTable
}

// RelationDatabase returns the endpoint type of the relation definition, which
// for a multi-table end lists every table it may point to, e.g. "square|circle".
func RelationDatabase(end EdgeEnd) string {
	names := make([]string, len(end.TargetTables()))

	for i, target := range end.TargetTables() {
		names[i] = target.NameDatabase()
	}

	return strings.Join(names, "|")
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

	// AssertConfig resolves a constraint referenced via `som:"assert=<name>"`
	// against the //go:build som definitions.
	AssertConfig func(name string) *parser.AssertDef
}

type baseField struct {
	*BuildConfig

	source parser.Field

	// lenFunc is the SurrealQL function measuring the length of the field's
	// database value, used by the len constraint.
	lenFunc string
}

// fieldDef holds the clauses of a DEFINE FIELD statement. Only Name, Table and
// Type are required, every other clause is omitted when empty.
type fieldDef struct {
	Name        string
	Table       string
	Type        string
	Default     string
	Value       string
	Assert      string
	Permissions string
}

// defineFieldStmt renders the clauses in the order SurrealDB expects them.
var defineFieldStmt = template.Must(template.New("defineField").Parse(
	`DEFINE FIELD OVERWRITE {{.Name}} ON TABLE {{.Table}} TYPE {{.Type}}` +
		`{{if .Default}} DEFAULT {{.Default}}{{end}}` +
		`{{if .Value}} VALUE {{.Value}}{{end}}` +
		`{{if .Assert}} ASSERT {{.Assert}}{{end}}` +
		`{{if .Permissions}} PERMISSIONS {{.Permissions}}{{end}};`,
))

// defineField renders the DEFINE FIELD statement for this field, merging the
// constraints declared via struct tags into the given built-in assert.
//
// The definition's Name is the full database path of the field, including the
// prefix of any struct it is nested in.
func (f *baseField) defineField(def fieldDef) string {
	def.Assert = f.assertExpression(def.Name, def.Type, def.Assert)

	var out strings.Builder
	if err := defineFieldStmt.Execute(&out, def); err != nil {
		// The template is a constant and every value is a plain string, so it
		// can only fail if the template itself is broken.
		panic(fmt.Sprintf("could not render field %s: %v", def.Name, err))
	}

	return out.String()
}

// define is the shorthand for the majority of fields, which define nothing
// beyond their name and database type.
func (f *baseField) define(table, prefix, dbType string) []string {
	return []string{
		f.defineField(fieldDef{
			Name:  prefix + f.NameDatabase(),
			Table: table,
			Type:  dbType,
		}),
	}
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

// IsPointer reports whether the model field is a pointer.
func (f *baseField) IsPointer() bool {
	return f.source.Pointer()
}

func (f *baseField) HasValidate() bool {
	return f.source.HasValidate()
}

func (f *baseField) SkipValidate() bool {
	return f.source.SkipValidate()
}
