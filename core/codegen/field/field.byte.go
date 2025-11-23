package field

import (
	"fmt"
	"math"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Byte struct {
	*baseField

	source *parser.FieldByte
}

func (f *Byte) typeGo() jen.Code {
	return jen.Add(f.ptr()).Byte()
}

func (f *Byte) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *Byte) TypeDatabase() string {
	return f.optionWrap("int")
}

func (f *Byte) SchemaStatements(table, prefix string) []string {
	nilCheck := ""
	if f.source.Pointer() {
		nilCheck = "$value == NONE OR $value == NULL OR "
	}

	extend := fmt.Sprintf("ASSERT %s$value >= %d AND $value <= %d", nilCheck, 0, math.MaxUint8)
	return []string{f.schemaStatement(table, prefix, f.TypeDatabase(), extend)}
}

func (f *Byte) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil, // Byte does not need a filter function.

		sortDefine: nil, // TODO: should be sortable
		sortInit:   nil,
		sortFunc:   nil, // Byte does not need a sort function.

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

// IMP: https://github.com/orgs/surrealdb/discussions/1451

func (f *Byte) filterDefine(ctx Context) jen.Code {
	filter := "Byte"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Byte) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewByte"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Byte) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Byte) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *Byte) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
