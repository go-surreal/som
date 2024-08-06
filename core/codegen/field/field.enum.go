package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
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

func (f *Enum) typeConv() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.model.NameGo()) // TODO: support other enum base types (atomic)
}

func (f *Enum) TypeDatabase() string {
	if !slices.Contains(f.values, "") {
		f.values = append(f.values, "") // TODO: add warning to output?!
	}

	sort.Strings(f.values)

	var values []string
	for _, value := range f.values {
		values = append(values, fmt.Sprintf(`"%s"`, value))
	}

	in := strings.Join(values, ", ")

	if f.source.Pointer() {
		return "option<string | null> ASSERT $value == NULL OR $value INSIDE [" + in + "]"
	}

	return "string ASSERT $value INSIDE [" + in + "]"
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
	filter := "Base"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, f.source.Typ))
}

func (f *Enum) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, f.source.Typ)),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Enum) convFrom(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Enum) convTo(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Enum) fieldDef(_ Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
