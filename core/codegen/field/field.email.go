package field

import (
	"fmt"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Email struct {
	*baseField

	source *parser.FieldEmail
}

func (f *Email) typeGo() jen.Code {
	// Email is defined in the generated som package, not the model package
	// Compute the som package path by going up from the model package
	// e.g., if model is "github.com/user/project/model", som is "github.com/user/project/gen/som"
	parts := strings.Split(f.SourcePkg, "/")
	if len(parts) > 0 {
		// Remove last segment (model package name) and add gen/som
		basePath := strings.Join(parts[:len(parts)-1], "/")
		somPkg := path.Join(basePath, "gen/som")
		return jen.Add(f.ptr()).Qual(somPkg, "Email")
	}
	// Fallback
	return jen.Add(f.ptr()).Qual(f.SourcePkg, "Email")
}

func (f *Email) typeConv(_ Context) jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *Email) TypeDatabase() string {
	if f.source.Pointer() {
		return "option<string>"
	}

	return "string"
}

func (f *Email) SchemaStatements(table, prefix string) []string {
	var extend string
	if f.source.Pointer() {
		extend = "ASSERT $value == NONE OR $value == NULL OR string::is::email($value)"
	} else {
		extend = `ASSERT $value == "" OR string::is::email($value)`
	}

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(), extend,
		),
	}
}

func (f *Email) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

func (f *Email) filterDefine(ctx Context) jen.Code {
	filter := "Email"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Email) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewEmail"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Email) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Email) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Email) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Email) cborMarshal(_ Context) jen.Code {
	convFuncName := "fromEmail"
	if f.source.Pointer() {
		convFuncName += "Ptr"
	}

	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id(convFuncName).Call(jen.Id("c").Dot(f.NameGo())),
		)
	}

	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id(convFuncName).Call(jen.Id("c").Dot(f.NameGo()))
}

func (f *Email) cborUnmarshal(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).BlockFunc(func(g *jen.Group) {
			g.Var().Id("convVal").Op("*").String()
			g.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
			g.Id("c").Dot(f.NameGo()).Op("=").Id("toEmailPtr").Call(jen.Id("convVal"))
		})
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("convVal").String()
		g.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
		g.Id("c").Dot(f.NameGo()).Op("=").Id("toEmail").Call(jen.Id("convVal"))
	})
}
