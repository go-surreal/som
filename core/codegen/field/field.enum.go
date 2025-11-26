package field

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Enum struct {
	*baseField

	source *parser.FieldEnum
	model  EnumModel
	values []string
}

func (f *Enum) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.model.NameGo())
}

func (f *Enum) typeConv(_ Context) jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.model.NameGo()) // TODO: support other enum base types (atomic)
}

func (f *Enum) TypeDatabase() string {
	if !slices.Contains(f.values, "") {
		f.values = append(f.values, "") // TODO: add warning to output?!
	}

	sort.Strings(f.values)

	var formattedValues []string
	for _, value := range f.values {
		formattedValues = append(formattedValues, fmt.Sprintf(`"%s"`, value))
	}

	literals := strings.Join(formattedValues, " | ")

	if f.source.Pointer() {
		return "option<" + literals + ">"
	}

	return literals
}

func (f *Enum) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *Enum) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *Enum) filterDefine(ctx Context) jen.Code {
	filter := "Enum"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel, jen.Qual(ctx.SourcePkg, f.source.Typ))
}

func (f *Enum) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewEnum"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel, jen.Qual(ctx.SourcePkg, f.source.Typ)),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Enum) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Enum) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Enum) cborMarshal(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *Enum) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
