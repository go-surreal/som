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

func (f *Time) typeConv(_ Context) jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgModels, "CustomDateTime")
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

func (f *Time) TypeDatabaseForArray() string {
	// Returns base type without VALUE/PERMISSIONS clauses for use in array element types
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
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Time) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewTime"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Time) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Time) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Time) convFrom(_ Context) (jen.Code, jen.Code) {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return nil, nil // never sent a timestamp to the database, as it will be set automatically
	}

	fromFunc := "fromTime"

	if f.source.Pointer() {
		fromFunc += fnSuffixPtr
	}

	return jen.Id(fromFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Time) convTo(ctx Context) (jen.Code, jen.Code) {
	if f.source.IsCreatedAt {
		return jen.Qual(ctx.TargetPkg, "NewTimestamps"),
			jen.Call(
				jen.Id("data").Dot("CreatedAt"),
				jen.Id("data").Dot("UpdatedAt"),
			)
	}

	if f.source.IsUpdatedAt {
		return nil, nil
	}

	toFunc := "toTime"

	if f.source.Pointer() {
		toFunc += fnSuffixPtr
	}

	return jen.Id(toFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Time) convToField(_ Context) jen.Code {
	if !f.source.IsCreatedAt {
		return nil
	}

	return jen.Id("Timestamps")
}

func (f *Time) fieldDef(ctx Context) jen.Code {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return jen.Id(f.NameGo()).Op("*").Add(f.typeConv(ctx)).
			Tag(map[string]string{convTag: f.NameDatabase() + ",omitempty"})
	}

	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase()})
}
