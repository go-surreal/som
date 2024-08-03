package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Time struct {
	*baseField

	source *parser.FieldTime
}

func (f *Time) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Time")
}

func (f *Time) typeConv() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgSDBC, "DateTime")
}

func (f *Time) TypeDatabase() string {
	if f.source.IsCreatedAt {
		// READONLY not working as expected, so using permissions as workaround for now.
		// See: https://surrealdb.com/docs/surrealdb/surrealql/statements/define/field#making-a-field-readonly-since-120
		return "option<datetime> VALUE $before OR time::now() PERMISSIONS FOR SELECT WHERE TRUE"
	}

	if f.source.IsUpdatedAt {
		// READONLY not working as expected, so using permissions as workaround for now.
		// See: https://surrealdb.com/docs/surrealdb/surrealql/statements/define/field#making-a-field-readonly-since-120
		return "option<datetime> VALUE time::now() PERMISSIONS FOR SELECT WHERE TRUE"
	}

	return f.optionWrap("datetime")
}

func (f *Time) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		convFrom:    f.convFrom,
		convTo:      f.convTo,
		convToField: f.convToField,
		fieldDef:    f.fieldDef,
	}
}

func (f *Time) filterDefine(ctx Context) jen.Code {
	filter := "Time"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Id("T"))
}

func (f *Time) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewTime"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Id("T")),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Time) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(jen.Id("T"))
}

func (f *Time) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Time) convFrom(_ Context) jen.Code {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return nil // never sent a timestamp to the database, as it will be set automatically
	}

	if f.source.Pointer() {
		return jen.Id("fromTimePtr").Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Qual(def.PkgSDBC, "DateTime").Values(jen.Id("data").Dot(f.NameGo()))
}

func (f *Time) convTo(_ Context) jen.Code {
	if f.source.IsCreatedAt {
		return jen.Qual(def.PkgSom, "NewTimestamps").Call(
			jen.Id("data").Dot("CreatedAt"),
			jen.Id("data").Dot("UpdatedAt"),
		)
	}

	if f.source.IsUpdatedAt {
		return nil
	}

	if f.source.Pointer() {
		return jen.Id("toTimePtr").Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Id("data").Dot(f.NameGo()).Dot("Time")
}

func (f *Time) convToField(_ Context) jen.Code {
	if !f.source.IsCreatedAt {
		return nil
	}

	return jen.Id("Timestamps")
}

func (f *Time) fieldDef(_ Context) jen.Code {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return jen.Id(f.NameGo()).Op("*").Add(f.typeConv()).
			Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
	}

	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
