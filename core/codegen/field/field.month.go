package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Month struct {
	*baseField

	source *parser.FieldMonth
}

func (f *Month) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Month")
}

func (f *Month) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *Month) TypeDatabase() string {
	return f.optionWrap("int")
}

func (f *Month) SchemaStatements(table, prefix string) []string {
	nilCheck := ""
	if f.source.Pointer() {
		nilCheck = "$value == NONE OR $value == NULL OR "
	}

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s ASSERT %s$value >= 1 AND $value <= 12;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(), nilCheck,
		),
	}
}

func (f *Month) CodeGen() *CodeGen {
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

func (f *Month) filterDefine(ctx Context) jen.Code {
	filter := "Month"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Month) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewMonth"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(ctx.filterKeyCode(f.NameDatabase()))
}

func (f *Month) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Month) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *Month) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgDistinct(), "Field").Types(def.TypeModel, jen.Qual("time", "Month"))
}

func (f *Month) fieldInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgDistinct(), "NewField").Types(def.TypeModel, jen.Qual("time", "Month")).
		Call(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *Month) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("val").Op(":=").Int().Call(jen.Op("*").Id("c").Dot(f.NameGo())),
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("val"),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Int().Call(jen.Id("c").Dot(f.NameGo()))
}

func (f *Month) cborUnmarshal(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).Block(
			jen.Var().Id("val").Op("*").Int(),
			jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("val")),
			jen.If(jen.Id("val").Op("!=").Nil()).Block(
				jen.Id("m").Op(":=").Qual("time", "Month").Call(jen.Op("*").Id("val")),
				jen.Id("c").Dot(f.NameGo()).Op("=").Op("&").Id("m"),
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
		jen.Id("c").Dot(f.NameGo()).Op("=").Qual("time", "Month").Call(jen.Id("val")),
	)
}

