package field

import (
	"fmt"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Duration struct {
	*baseField

	source *parser.FieldDuration
}

func (f *Duration) TypeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Duration")
}

func (f *Duration) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "Duration")
}

func (f *Duration) TypeDatabase() string {
	return f.optionWrap("duration")
}

func (f *Duration) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *Duration) CodeGen() *CodeGen {
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

		selectDecode:     f.selectDecode,
		selectDistDecode: f.selectDistDecode,
	}
}

func (f *Duration) selectDecode(ctx Context) jen.Code {
	dt := jen.Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "Duration")

	if f.source.Pointer() {
		return jen.Func().Params(jen.Id("data").Index().Byte()).Params(jen.Index().Op("*").Qual("time", "Duration"), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalSelectConvert").Call(jen.Id("data"), jen.Func().Params(jen.Id("v").Op("*").Add(dt)).Op("*").Qual("time", "Duration").Block(
				jen.If(jen.Id("v").Op("==").Nil()).Block(jen.Return(jen.Nil())),
				jen.Id("d").Op(":=").Id("v").Dot("Duration"),
				jen.Return(jen.Op("&").Id("d")),
			))))
	}

	return jen.Func().Params(jen.Id("data").Index().Byte()).Params(jen.Index().Qual("time", "Duration"), jen.Error()).Block(
		jen.Return(jen.Id("unmarshalSelectConvert").Call(jen.Id("data"), jen.Func().Params(jen.Id("v").Add(dt)).Qual("time", "Duration").Block(
			jen.Return(jen.Id("v").Dot("Duration")),
		))))
}

func (f *Duration) selectDistDecode(ctx Context) jen.Code {
	dt := jen.Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "Duration")

	if f.source.Pointer() {
		return jen.Func().Params(jen.Id("data").Index().Byte()).Params(jen.Index().Op("*").Qual("time", "Duration"), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalSelectDistinctConvert").Call(jen.Id("data"), jen.Func().Params(jen.Id("v").Op("*").Add(dt)).Op("*").Qual("time", "Duration").Block(
				jen.If(jen.Id("v").Op("==").Nil()).Block(jen.Return(jen.Nil())),
				jen.Id("d").Op(":=").Id("v").Dot("Duration"),
				jen.Return(jen.Op("&").Id("d")),
			))))
	}

	return jen.Func().Params(jen.Id("data").Index().Byte()).Params(jen.Index().Qual("time", "Duration"), jen.Error()).Block(
		jen.Return(jen.Id("unmarshalSelectDistinctConvert").Call(jen.Id("data"), jen.Func().Params(jen.Id("v").Add(dt)).Qual("time", "Duration").Block(
			jen.Return(jen.Id("v").Dot("Duration")),
		))))
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
		jen.Params(ctx.filterKeyCode(f.NameDatabase()))
}

func (f *Duration) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Duration) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *Duration) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgDistinct(), "Field").Types(def.TypeModel, jen.Qual("time", "Duration"))
}

func (f *Duration) fieldInit(ctx Context) jen.Code {
	factory := "NewDurationField"
	if f.source.Pointer() {
		factory = "NewDurationPtrField"
	}
	return jen.Qual(ctx.pkgDistinct(), factory).Types(def.TypeModel).
		Call(ctx.sortKeyCode(f.NameDatabase()))
}

func (f *Duration) cborMarshal(ctx Context) jen.Code {
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
		jen.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(ctx.pkgCBOR(), helper).Call(jen.Id("raw")),
	)
}
