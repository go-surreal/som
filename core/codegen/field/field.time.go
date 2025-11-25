package field

import (
	"fmt"
	"path"

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

func (f *Time) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "DateTime")
}

func (f *Time) TypeDatabase() string {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return "option<datetime>"
	}

	return f.optionWrap("datetime")
}

func (f *Time) SchemaStatements(table, prefix string) []string {
	var extend string

	if f.source.IsCreatedAt {
		extend = "VALUE $before OR time::now() READONLY"
	} else if f.source.IsUpdatedAt {
		extend = "VALUE time::now()"
	}

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(), extend,
		),
	}
}

func (f *Time) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

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

func (f *Time) fieldDef(ctx Context) jen.Code {
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return jen.Id(f.NameGo()).Op("*").Add(f.typeConv(ctx)).
			Tag(map[string]string{convTag: f.NameDatabase() + ",omitempty"})
	}

	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Time) cborMarshal(ctx Context) jen.Code {
	// Timestamp fields use getter methods from embedded Timestamps.
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		return jen.If(jen.Op("!").Id("c").Dot(f.NameGo()).Call().Dot("IsZero").Call()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "DateTime").Values(
				jen.Id("Time").Op(":").Id("c").Dot(f.NameGo()).Call(),
			),
		)
	}

	// Using custom types.DateTime with MarshalCBOR method.
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "DateTime").Values(
				jen.Id("Time").Op(":").Op("*").Id("c").Dot(f.NameGo()),
			),
		)
	}

	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "DateTime").Values(
		jen.Id("Time").Op(":").Id("c").Dot(f.NameGo()),
	)
}

func (f *Time) cborUnmarshal(ctx Context) jen.Code {
	// Timestamp fields use setter methods on embedded Timestamps.
	if f.source.IsCreatedAt || f.source.IsUpdatedAt {
		setter := "SetCreatedAt"
		if f.source.IsUpdatedAt {
			setter = "SetUpdatedAt"
		}
		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).Block(
			jen.Id("tm").Op(",").Id("_").Op(":=").Qual(ctx.pkgCBOR(), "UnmarshalDateTime").Call(jen.Id("raw")),
			jen.Id("c").Dot("Timestamps").Dot(setter).Call(jen.Id("tm")),
		)
	}

	helper := "UnmarshalDateTime"
	if f.source.Pointer() {
		helper = "UnmarshalDateTimePtr"
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(ctx.pkgCBOR(), helper).Call(jen.Id("raw")),
	)
}
