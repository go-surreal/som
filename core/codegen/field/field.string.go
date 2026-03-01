package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type String struct {
	*baseField

	source *parser.FieldString
}

func (f *String) typeGo() jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *String) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *String) TypeDatabase() string {
	return f.optionWrap("string")
}

func (f *String) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *String) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,
		filterExtra:  f.filterExtra,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		fieldDefine: f.fieldDefine,
		fieldInit:   f.fieldInit,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *String) filterDefine(ctx Context) jen.Code {
	// For searchable (index) strings, we use a wrapper type (see filterExtra below).
	if f.SearchInfo() != nil {
		return jen.Id(f.NameGo()).Id(ctx.Table.NameGoLower() + f.NameGo()).Types(def.TypeModel)
	}

	filter := "String"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *String) filterInit(ctx Context) (jen.Code, jen.Code) {
	// For searchable (index) strings, we use a wrapper type (see filterExtra below).
	if f.SearchInfo() != nil {
		wrapperName := ctx.Table.NameGoLower() + f.NameGo()
		filter := "NewString"

		if f.source.Pointer() {
			filter += fnSuffixPtr
		}

		return jen.Id(wrapperName).Types(def.TypeModel).Values(
			jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel).
				Call(ctx.filterKeyCode(f.NameDatabase())),
		), jen.Empty()
	}

	filter := "NewString"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(ctx.filterKeyCode(f.NameDatabase()))
}

// filterExtra generates the wrapper type and Matches() method for search-indexed strings.
func (f *String) filterExtra(ctx Context) jen.Code {
	if f.SearchInfo() == nil {
		return nil
	}

	wrapperName := ctx.Table.NameGoLower() + f.NameGo()
	embeddedType := "String"
	if f.source.Pointer() {
		embeddedType += fnSuffixPtr
	}

	return jen.Add(
		jen.Type().Id(wrapperName).Types(jen.Add(def.TypeModel).Any()).Struct(
			jen.Op("*").Qual(ctx.pkgLib(), embeddedType).Types(def.TypeModel),
		),
		jen.Line(),
		jen.Func().
			Params(jen.Id("f").Id(wrapperName).Types(def.TypeModel)).
			Id("Matches").
			Params(jen.Id("terms").String()).
			Qual(ctx.pkgLib(), "Search").Types(def.TypeModel).
			Block(
				jen.Return(
					jen.Qual(ctx.pkgLib(), "NewSearch").Types(def.TypeModel).Call(
						jen.Id("f").Dot(embeddedType).Dot("Base").Dot("Key"),
						jen.Id("terms"),
					),
				),
			),
		jen.Line(),
		jen.Func().
			Params(jen.Id("f").Id(wrapperName).Types(def.TypeModel)).
			Id("MatchesAny").
			Params(jen.Id("terms").String()).
			Qual(ctx.pkgLib(), "Search").Types(def.TypeModel).
			Block(
				jen.Return(
					jen.Qual(ctx.pkgLib(), "NewSearchOr").Types(def.TypeModel).Call(
						jen.Id("f").Dot(embeddedType).Dot("Base").Dot("Key"),
						jen.Id("terms"),
					),
				),
			),
		jen.Line(),
		jen.Func().
			Params(jen.Id("f").Id(wrapperName).Types(def.TypeModel)).
			Id("key").
			Params().
			Qual(ctx.pkgLib(), "Key").Types(def.TypeModel).
			Block(
				jen.Return(jen.Id("f").Dot(embeddedType).Dot("Base").Dot("Key")),
			),
	)
}

func (f *String) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "StringSort").Types(def.TypeModel)
}

func (f *String) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewStringSort").Types(def.TypeModel).
		Params(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *String) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgDistinct(), "Field").Types(def.TypeModel, jen.String())
}

func (f *String) fieldInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgDistinct(), "NewField").Types(def.TypeModel, jen.String()).
		Call(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *String) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *String) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
