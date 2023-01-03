package field

import (
	"encoding/json"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
	"golang.org/x/exp/slices"
	"sort"
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
	return jen.Add(f.ptr()).String() // TODO: support other enum base types (atomic)
}

func (f *Enum) TypeDatabase() string {
	if !slices.Contains(f.values, "") {
		f.values = append(f.values, "") // TODO: add warning to output?!
	}

	sort.Strings(f.values)
	valuesRaw, _ := json.Marshal(f.values)
	values := string(valuesRaw)

	if f.source.Pointer() {
		return "string ASSERT $value == NULL OR $value INSIDE " + values
	}

	return "string ASSERT $value INSIDE " + values
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

	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, filter).Types(jen.Qual(ctx.SourcePkg, f.source.Typ), jen.Id("T"))
}

func (f *Enum) filterInit(ctx Context) jen.Code {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(def.PkgLibFilter, filter).Types(jen.Qual(ctx.SourcePkg, f.source.Typ), jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(f.NameDatabase())))
}

func (f *Enum) convFrom(ctx Context) jen.Code {
	return jen.String().Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Enum) convTo(ctx Context) jen.Code {
	return jen.Qual(ctx.SourcePkg, f.source.Typ).Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Enum) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
