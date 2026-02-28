package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Weekday struct {
	*baseField

	source *parser.FieldWeekday
}

func (f *Weekday) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Weekday")
}

func (f *Weekday) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *Weekday) TypeDatabase() string {
	return f.optionWrap("int")
}

func (f *Weekday) SchemaStatements(table, prefix string) []string {
	nilCheck := ""
	if f.source.Pointer() {
		nilCheck = "$value == NONE OR $value == NULL OR "
	}

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s ASSERT %s$value >= 0 AND $value <= 6;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(), nilCheck,
		),
	}
}

func (f *Weekday) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		fieldDefine: f.fieldDefine,
		fieldInit:   f.fieldInit,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *Weekday) filterDefine(ctx Context) jen.Code {
	filter := "Weekday"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Weekday) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewWeekday"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(ctx.filterKeyCode(f.NameDatabase()))
}

func (f *Weekday) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Weekday) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *Weekday) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgDistinct(), "Field").Types(def.TypeModel, jen.Qual("time", "Weekday"))
}

func (f *Weekday) fieldInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgDistinct(), "NewField").Types(def.TypeModel, jen.Qual("time", "Weekday")).
		Call(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *Weekday) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("val").Op(":=").Int().Call(jen.Op("*").Id("c").Dot(f.NameGo())),
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("val"),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Int().Call(jen.Id("c").Dot(f.NameGo()))
}

func (f *Weekday) cborUnmarshal(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).Block(
			jen.Var().Id("val").Op("*").Int(),
			jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("val")),
			jen.If(jen.Id("val").Op("!=").Nil()).Block(
				jen.Id("w").Op(":=").Qual("time", "Weekday").Call(jen.Op("*").Id("val")),
				jen.Id("c").Dot(f.NameGo()).Op("=").Op("&").Id("w"),
			).Else().Block(
				jen.Id("c").Dot(f.NameGo()).Op("=").Nil(),
			),
		)
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Var().Id("val").Int(),
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("val")),
		jen.Id("c").Dot(f.NameGo()).Op("=").Qual("time", "Weekday").Call(jen.Id("val")),
	)
}

