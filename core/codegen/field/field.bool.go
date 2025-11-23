package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Bool struct {
	*baseField

	source *parser.FieldBool
}

func (f *Bool) typeGo() jen.Code {
	return jen.Add(f.ptr()).Bool()
}

func (f *Bool) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *Bool) TypeDatabase() string {
	return f.optionWrap("bool")
}

func (f *Bool) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *Bool) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil, // Bool does not need a filter function.

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

func (f *Bool) filterDefine(ctx Context) jen.Code {
	filter := "Bool"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Bool) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewBool"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Bool) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Bool) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Bool) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Bool) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *Bool) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
