package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
	"slices"
	"sort"
	"strings"
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
		return "option<" + literals + " | null>"
	}

	return literals
}

func (f *Enum) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil,

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
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

func (f *Enum) convFrom(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Enum) convTo(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Enum) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase()})
}
