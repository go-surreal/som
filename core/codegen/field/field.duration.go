package field

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Duration struct {
	*baseField

	source *parser.FieldDuration
}

func (f *Duration) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Duration")
}

func (f *Duration) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "Duration")
}

func (f *Duration) TypeDatabase() string {
	return f.optionWrap("duration")
}

func (f *Duration) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		convFrom:      f.convFrom,
		convTo:        f.convTo,
		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

func (f *Duration) filterDefine(ctx Context) jen.Code {
	filter := "Duration"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Duration) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewDuration"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Duration) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Duration) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Duration) convFrom(_ Context) (jen.Code, jen.Code) {
	fromFunc := "fromDuration"

	if f.source.Pointer() {
		fromFunc += fnSuffixPtr
	}

	return jen.Id(fromFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Duration) convTo(_ Context) (jen.Code, jen.Code) {
	toFunc := "toDuration"

	if f.source.Pointer() {
		toFunc += fnSuffixPtr
	}

	return jen.Id(toFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Duration) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Duration) cborMarshal(ctx Context) jen.Code {
	// Direct assignment - types.Duration has MarshalCBOR that cbor.Marshal will call
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "Duration").Values(
				jen.Id("Duration").Op(":").Op("*").Id("c").Dot(f.NameGo()),
			),
		)
	}

	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "Duration").Values(
		jen.Id("Duration").Op(":").Id("c").Dot(f.NameGo()),
	)
}

func (f *Duration) cborUnmarshal(ctx Context) jen.Code {
	helper := "UnmarshalDuration"
	if f.source.Pointer() {
		helper = "UnmarshalDurationPtr"
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(ctx.pkgCBORHelpers(), helper).Call(jen.Id("raw")),
	)
}
