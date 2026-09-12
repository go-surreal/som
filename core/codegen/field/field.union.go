package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

// Union is a record link that may point to any member of a union, which maps
// to a multi-table record type (record<a|b>) in the database.
type Union struct {
	*baseField

	source *parser.FieldUnion
	union  *UnionTable
}

func (f *Union) typeGo() jen.Code {
	return jen.Qual(f.SourcePkg, f.union.NameGo())
}

func (f *Union) typeConv(_ Context) jen.Code {
	return jen.Op("*").Id(f.union.NameGoLower() + "Link")
}

func (f *Union) TypeDatabase() string {
	// Linked records are always considered optional.
	return fmt.Sprintf("option<%s>", f.union.TypeDatabase())
}

func (f *Union) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD OVERWRITE %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *Union) Union() *UnionTable {
	return f.union
}

func (f *Union) TargetNameGo() string {
	return f.union.NameGo()
}

func (f *Union) TargetNameGoLower() string {
	return f.union.NameGoLower()
}

func (f *Union) TargetTables() []*NodeTable {
	return f.union.Members
}

func (f *Union) RelationDatabase() string {
	return f.union.RelationDatabase()
}

func (f *Union) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   f.filterFunc,

		fieldFunc: f.fieldFieldFunc,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *Union) filterDefine(_ Context) jen.Code {
	return jen.Id(f.union.NameGoLower()).Types(def.TypeModel)
}

func (f *Union) filterInit(_ Context) (jen.Code, jen.Code) {
	return jen.Id("new" + f.union.NameGo()).Types(def.TypeModel), nil
}

func (f *Union) filterFunc(ctx Context) jen.Code {
	receiver := jen.Id(ctx.Table.NameGoLower()).Types(def.TypeModel)
	if ctx.Receiver != nil {
		receiver = ctx.Receiver
	}

	return jen.Func().
		Params(jen.Id("n").Add(receiver)).
		Id(f.NameGo()).Params().
		Add(f.filterDefine(ctx)).
		Block(
			jen.Return(jen.Add(f.filterInit(ctx)).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())))))
}

func (f *Union) fieldFieldFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Id(f.union.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id("new" + f.union.NameGo()).Types(def.TypeModel).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Union) cborMarshal(_ Context) jen.Code {
	// A union field holds an interface, so the concrete member it points to
	// determines the table of the record link.
	return jen.If(
		jen.Id("link").Op(":=").Id("to"+f.union.NameGo()+"Link").Call(jen.Id("c").Dot(f.NameGo())),
		jen.Id("link").Op("!=").Nil(),
	).Block(
		jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("link"),
	)
}

func (f *Union) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("convVal").Op("*").Id(f.union.NameGoLower() + "Link")
		g.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
		g.Id("c").Dot(f.NameGo()).Op("=").Id("from" + f.union.NameGo() + "Link").Call(jen.Id("convVal"))
	})
}
