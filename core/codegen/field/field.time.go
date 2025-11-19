package field

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Time struct {
	*baseField

	source *parser.FieldTime
}

func (f *Time) Source() *parser.FieldTime {
	return f.source
}

func (f *Time) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Time")
}

func (f *Time) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "DateTime")
}

func (f *Time) TypeDatabase() string {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return "option<datetime>"
	}

	return f.optionWrap("datetime")
}

func (f *Time) TypeDatabaseExtend() string {
	if f.source.IsCreatedAt {
		// READONLY not working as expected, so using permissions as workaround for now.
		// See: https://surrealdb.com/docs/surrealdb/surrealql/statements/define/field#making-a-field-readonly-since-120
		return "VALUE $before OR time::now() PERMISSIONS FOR SELECT WHERE TRUE"
	}

	if f.source.IsUpdatedAt {
		// READONLY not working as expected, so using permissions as workaround for now.
		// See: https://surrealdb.com/docs/surrealdb/surrealql/statements/define/field#making-a-field-readonly-since-120
		return "VALUE time::now() PERMISSIONS FOR SELECT WHERE TRUE"
	}

	return ""
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

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,

		fieldDef: f.fieldDef,
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
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Time) cborMarshal(ctx Context) jen.Code {
	// Skip timestamp fields - they're handled separately in buildMarshalCBOR
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return nil
	}

	helper := "marshalDateTime"
	if f.source.Pointer() {
		helper += "Ptr"
	}

	return jen.BlockFunc(func(g *jen.Group) {
		if f.source.Pointer() {
			// For pointers, check if non-nil before marshaling
			g.If(jen.Id("m").Dot(f.NameGo()).Op("!=").Nil()).Block(
				jen.Id("val").Op(",").Id("_").Op(":=").Qual(path.Join(ctx.TargetPkg, def.PkgCBOR), helper).Call(
					jen.Id("m").Dot(f.NameGo()),
				),
				jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Qual(def.PkgCBOR, "RawMessage").Call(jen.Id("val")),
			)
		} else {
			// For non-pointers, always marshal
			g.Id("val").Op(",").Id("_").Op(":=").Qual(path.Join(ctx.TargetPkg, def.PkgCBOR), helper).Call(
				jen.Id("m").Dot(f.NameGo()),
			)
			g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Qual(def.PkgCBOR, "RawMessage").Call(jen.Id("val"))
		}
	})
}

func (f *Time) cborUnmarshal(ctx Context) jen.Code {
	// Skip timestamp fields - they're handled separately in buildUnmarshalCBOR
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return nil
	}

	helper := "unmarshalDateTime"
	if f.source.Pointer() {
		helper += "Ptr"
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Id("m").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(path.Join(ctx.TargetPkg, def.PkgCBOR), helper).Call(jen.Id("raw")),
	)
}
