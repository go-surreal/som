package field

import (
	"fmt"
	"slices"
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

// literals renders the declared enum values as a SurrealDB literal union.
func (f *Enum) literals() string {
	formattedValues := make([]string, 0, len(f.values))
	for _, value := range f.values {
		formattedValues = append(formattedValues, fmt.Sprintf(`"%s"`, value))
	}

	return strings.Join(formattedValues, " | ")
}

// hasEmptyValue reports whether the empty string is a declared enum value. Only
// then is it a variant in its own right, instead of merely the Go zero value.
func (f *Enum) hasEmptyValue() bool {
	return slices.Contains(f.values, "")
}

func (f *Enum) TypeDatabase() string {
	// A non-pointer enum has no way to express "unset" in Go other than the zero
	// value, which is stored as NONE unless the empty string is a declared value.
	if f.source.Pointer() || !f.hasEmptyValue() {
		return "option<" + f.literals() + ">"
	}

	return f.literals()
}

// ElementTypeDatabase returns the type for the enum as a slice element. Elements
// are never implicitly zero, so the Go zero value needs no representation here.
func (f *Enum) ElementTypeDatabase() string {
	return f.optionWrap(f.literals())
}

func (f *Enum) SchemaStatements(table, prefix string) []string {
	return f.define(table, prefix, f.TypeDatabase())
}

func (f *Enum) CodeGen() *CodeGen {
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

func (f *Enum) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgDistinct(), "Field").Types(def.TypeModel, jen.Qual(ctx.SourcePkg, f.source.Typ))
}

func (f *Enum) fieldInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgDistinct(), "NewField").Types(def.TypeModel, jen.Qual(ctx.SourcePkg, f.source.Typ)).
		Call(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Enum) cborMarshal(_ Context) jen.Code {
	assign := jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())

	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(assign)
	}

	if !f.hasEmptyValue() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Lit("")).Block(assign)
	}

	return assign
}

func (f *Enum) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
