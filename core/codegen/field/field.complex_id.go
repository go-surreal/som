package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type ComplexID struct {
	*baseField

	source  *parser.FieldComplexID
	element Table
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
		typeDef = fmt.Sprintf("array<any, %d>", len(f.source.Fields))
	case parser.IDTypeObject:
		return []string{
			fmt.Sprintf("DEFINE FIELD id ON TABLE %s TYPE object FLEXIBLE;", table),
		}
	default:
		return nil
	}
	return []string{
		fmt.Sprintf("DEFINE FIELD id ON TABLE %s TYPE %s;", table, typeDef),
	}
}

func (f *ComplexID) CodeGen() *CodeGen {
	if f.element == nil {
		return &CodeGen{}
	}
	return &CodeGen{
		filterFunc: f.filterFunc,
		sortFunc:   f.sortFunc,
		fieldFunc:  f.fieldFieldFunc,
	}
}

func (f *ComplexID) filterFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Id(f.element.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id("new"+f.source.StructName).Types(def.TypeModel).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())))))
}

func (f *ComplexID) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Id(f.element.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id("new"+f.source.StructName).Types(def.TypeModel).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *ComplexID) fieldFieldFunc(ctx Context) jen.Code {
	return f.sortFunc(ctx)
}
