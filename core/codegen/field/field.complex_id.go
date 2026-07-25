package field

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

// idPart holds the database name and type of a single complex ID field,
// used to build a literal (tuple/object) type for the record ID.
type idPart struct {
	dbName string
	dbType string
}

type ComplexID struct {
	*baseField

	source  *parser.FieldComplexID
	element Table
	parts   []idPart
}

func (f *ComplexID) typeGo() jen.Code {
	return jen.Qual(f.SourcePkg, f.source.StructName)
}

func (f *ComplexID) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *ComplexID) TypeDatabase() string {
	switch f.source.Kind {
	case parser.IDTypeArray:
		return "array"
	case parser.IDTypeObject:
		return "object"
	default:
		return "any"
	}
}

func (f *ComplexID) SchemaStatements(table, _ string) []string {
	var typeDef string
	switch f.source.Kind {
	case parser.IDTypeArray:
		types := make([]string, len(f.parts))
		for i, p := range f.parts {
			types[i] = p.dbType
		}
		typeDef = "[" + strings.Join(types, ", ") + "]"
	case parser.IDTypeObject:
		entries := make([]string, len(f.parts))
		for i, p := range f.parts {
			entries[i] = fmt.Sprintf("%s: %s", p.dbName, p.dbType)
		}
		typeDef = "{ " + strings.Join(entries, ", ") + " }"
	default:
		return nil
	}
	return []string{
		fmt.Sprintf("DEFINE FIELD id ON TABLE %s TYPE %s;", table, typeDef),
	}
}

// idPartType returns the database type of a complex ID field. ID parts are
// always required, so any outer option<...> wrapper is stripped.
func idPartType(f Field) string {
	t := f.TypeDatabase()
	if strings.HasPrefix(t, "option<") && strings.HasSuffix(t, ">") {
		t = t[len("option<") : len(t)-1]
	}
	return t
}

func (f *ComplexID) CodeGen() *CodeGen {
	if f.element == nil {
		return &CodeGen{}
	}
	return &CodeGen{
		filterFunc: f.filterFunc,
	}
}

func (f *ComplexID) filterFunc(ctx Context) jen.Code {
	idKey := jen.Qual(ctx.pkgLib(), "Fn").Call(
		jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())),
		jen.Lit("meta::id"),
	)

	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Id(f.element.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id("new"+f.source.StructName).Types(def.TypeModel).
				Params(idKey)))
}
