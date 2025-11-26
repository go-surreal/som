package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type ID struct {
	*baseField

	source *parser.FieldID
}

func (f *ID) typeGo() jen.Code {
	return jen.String()
}

func (f *ID) typeConv(ctx Context) jen.Code {
	return jen.Op("*").Qual(ctx.TargetPkg, "ID") // f.typeGo()
}

func (f *ID) TypeDatabase() string {
	// TODO: type "uuid" works, but there is no native type "ulid"
	// see: https://github.com/surrealdb/surrealdb/issues/1722
	return "string"
}

func (f *ID) SchemaStatements(table, prefix string) []string {
	// TODO: assert := "string::is::ulid(record::id($value))"

	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *ID) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,
	}
}

func (f *ID) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "ID").Types(def.TypeModel)
}

func (f *ID) filterInit(ctx Context) (jen.Code, jen.Code) {
	return jen.Qual(ctx.pkgLib(), "NewID").Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())), jen.Lit(ctx.Table.NameDatabase()))
}

func (f *ID) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *ID) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}
