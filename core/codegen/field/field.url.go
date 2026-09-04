package field

import (
	"fmt"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type URL struct {
	*baseField

	source *parser.FieldURL
}

func (f *URL) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgURL, "URL")
}

func (f *URL) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *URL) TypeDatabase() string {
	if f.source.Pointer() {
		return "option<string>"
	}

	return "string"
}

func (f *URL) SchemaStatements(table, prefix string) []string {
	var extend string
	if f.source.Pointer() {
		extend = "ASSERT $value == NONE OR $value == NULL OR string::is_url($value)"
	} else {
		extend = `ASSERT $value == "" OR string::is_url($value)`
	}

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD OVERWRITE %s ON TABLE %s TYPE %s %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(), extend,
		),
	}
}

func (f *URL) CodeGen() *CodeGen {
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

func (f *URL) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgDistinct(), "Field").Types(def.TypeModel, jen.Qual(def.PkgURL, "URL"))
}

func (f *URL) fieldInit(ctx Context) jen.Code {
	factory := "NewURLField"
	if f.source.Pointer() {
		factory = "NewURLPtrField"
	}
	return jen.Qual(ctx.pkgDistinct(), factory).Types(def.TypeModel).
		Call(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *URL) filterDefine(ctx Context) jen.Code {
	filter := "URL"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *URL) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewURL"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *URL) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *URL) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *URL) cborMarshal(ctx Context) jen.Code {
	typeURL := jen.Op("*").Qual(path.Join(ctx.TargetPkg, def.PkgTypes), "URL")

	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").
				Params(typeURL).Call(jen.Id("c").Dot(f.NameGo())),
		)
	}

	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").
		Params(typeURL).Call(jen.Op("&").Id("c").Dot(f.NameGo()))
}

func (f *URL) cborUnmarshal(ctx Context) jen.Code {
	typesPkg := path.Join(ctx.TargetPkg, def.PkgTypes)

	if f.source.Pointer() {
		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).BlockFunc(func(g *jen.Group) {
			g.If(jen.Qual(ctx.pkgCBOR(), "IsNoneOrNull").Call(jen.Id("raw"))).Block(
				jen.Id("c").Dot(f.NameGo()).Op("=").Nil(),
			).Else().Block(
				jen.Var().Id("convVal").Qual(def.PkgURL, "URL"),
				jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(
					jen.Id("raw"),
					jen.Params(jen.Op("*").Qual(typesPkg, "URL")).Call(jen.Op("&").Id("convVal")),
				),
				jen.Id("c").Dot(f.NameGo()).Op("=").Op("&").Id("convVal"),
			)
		})
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(
			jen.Id("raw"),
			jen.Params(jen.Op("*").Qual(typesPkg, "URL")).Call(jen.Op("&").Id("c").Dot(f.NameGo())),
		),
	)
}
