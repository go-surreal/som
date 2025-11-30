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

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *String) filterDefine(ctx Context) jen.Code {
	// For search-indexed strings, use a wrapper type
	if f.SearchInfo() != nil {
		// Returns: FieldName tableNameFieldName[M]
		wrapperName := ctx.Table.NameGoLower() + f.NameGo()
		return jen.Id(f.NameGo()).Id(wrapperName).Types(def.TypeModel)
	}

	// For non-indexed strings, use lib.String or lib.StringPtr
	filter := "String"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *String) filterInit(ctx Context) (jen.Code, jen.Code) {
	// For search-indexed strings, use wrapper type initialization
	if f.SearchInfo() != nil {
		wrapperName := ctx.Table.NameGoLower() + f.NameGo()
		filter := "NewString"
		if f.source.Pointer() {
			filter += fnSuffixPtr
		}
		// Returns: tableNameFieldName[M]{lib.NewString[M](lib.Field(key, "field_name"))}
		return jen.Id(wrapperName).Types(def.TypeModel).Values(
			jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel).
				Call(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase()))),
		), jen.Empty()
	}

	// For non-indexed strings
	filter := "NewString"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
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

	// Generate:
	// type tableNameFieldName[M any] struct {
	//     *lib.String[M]  // or *lib.StringPtr[M]
	// }
	//
	// func (f tableNameFieldName[M]) Matches(terms string) lib.Search[M] {
	//     return lib.NewSearch[M](f.String.Base.Key, terms)
	// }
	//
	// func (f tableNameFieldName[M]) key() lib.Key[M] {
	//     return f.String.Base.Key
	// }
	return jen.Add(
		// Wrapper type definition
		jen.Type().Id(wrapperName).Types(jen.Add(def.TypeModel).Any()).Struct(
			jen.Op("*").Qual(ctx.pkgLib(), embeddedType).Types(def.TypeModel),
		),
		jen.Line(),
		// Matches() method
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
		// key() method - implements field[M] interface for use with Equal_() etc.
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
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
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
