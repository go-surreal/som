package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

const versionDBField = "__som_lock_version"

type Version struct {
	*baseField

	source *parser.FieldVersion
}

func (f *Version) NameDatabase() string {
	return versionDBField
}

func (f *Version) typeGo() jen.Code {
	return jen.Int()
}

func (f *Version) typeConv(_ Context) jen.Code {
	return jen.Int()
}

func (f *Version) TypeDatabase() string {
	return "int"
}

func (f *Version) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			`DEFINE FIELD %s ON TABLE %s TYPE %s VALUE { IF $value != NONE AND $before != NONE AND $value != $before { THROW "optimistic_lock_failed" }; RETURN IF $before THEN $before + 1 ELSE 1 END; };`,
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *Version) CodeGen() *CodeGen {
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

func (f *Version) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "Int").Types(def.TypeModel, jen.Int())
}

func (f *Version) filterInit(ctx Context) (jen.Code, jen.Code) {
	return jen.Qual(ctx.pkgLib(), "NewInt").Types(def.TypeModel, jen.Int()),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Version) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *Version) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Version) fieldDef(_ Context) jen.Code {
	return jen.Id(f.NameGo()).Int().
		Tag(map[string]string{convTag: f.NameDatabase()})
}

func (f *Version) cborMarshal(_ Context) jen.Code {
	// Version() is a getter method on the embedded OptimisticLock struct
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()).Call()
}

func (f *Version) cborUnmarshal(ctx Context) jen.Code {
	// SetVersion is a setter method on the embedded OptimisticLock struct
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Var().Id("v").Int(),
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("v")),
		jen.Id("c").Dot("OptimisticLock").Dot("SetVersion").Call(jen.Id("v")),
	)
}
