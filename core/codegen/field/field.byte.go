package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
	"math"
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
	nilCheck := ""
	if f.source.Pointer() {
		nilCheck = "$value == NONE OR $value == NULL OR "
	}

	return fmt.Sprintf(
		"%s ASSERT %s$value >= %d AND $value <= %d",
		f.optionWrap("int"), nilCheck, 0, math.MaxUint8,
	)
}

func (f *Byte) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil, // Byte does not need a filter function.

		sortDefine: nil, // TODO: should be sortable
		sortInit:   nil,
		sortFunc:   nil, // Byte does not need a sort function.

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

// IMP: https://github.com/orgs/surrealdb/discussions/1451

func (f *Byte) filterDefine(ctx Context) jen.Code {
	filter := "Base"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Id("T"), jen.Byte())
}

func (f *Byte) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Id("T"), jen.Byte()),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Byte) convFrom(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Byte) convTo(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Byte) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{"json": f.NameDatabase()})
}
