package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Struct struct {
	*baseField

	source *parser.FieldStruct
	table  Table
}

func (f *Struct) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.table.NameGo())
}

func (f *Struct) typeConv(_ Context) jen.Code {
	return jen.Add(f.ptr()).Id(f.table.NameGoLower())
}

func (f *Struct) TypeDatabase() string {
	return f.optionWrap("object")
}

func (f *Struct) SchemaStatements(table, prefix string) []string {
	var statements []string

	// Generate own DEFINE FIELD statement
	statement := f.schemaStatement(table, prefix, f.TypeDatabase(), "")
	statements = append(statements, statement)

	// Recursively get nested field statements
	nestedPrefix := prefix + f.NameDatabase() + "."
	for _, fld := range f.table.GetFields() {
		statements = append(statements, fld.SchemaStatements(table, nestedPrefix)...)
	}

	return statements
}

func (f *Struct) Table() Table {
	return f.table
}

func (f *Struct) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil, // TODO

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

func (f *Struct) filterDefine(_ Context) jen.Code {
	return jen.Id(f.table.NameGoLower()).Types(def.TypeModel)
}

func (f *Struct) filterInit(_ Context) (jen.Code, jen.Code) {
	return jen.Id("new" + f.source.Struct).Types(def.TypeModel), nil
}

func (f *Struct) filterFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Add(f.filterDefine(ctx)).
		Block(
			jen.Return(jen.Add(f.filterInit(ctx)).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())))))
}

func (f *Struct) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Struct) cborMarshal(_ Context) jen.Code {
	// Convert to conv wrapper which has proper MarshalCBOR
	convFuncName := "from" + f.table.NameGo()
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

func (f *Struct) cborUnmarshal(ctx Context) jen.Code {
	// Unmarshal through conv wrapper
	if f.source.Pointer() {
		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).BlockFunc(func(g *jen.Group) {
			g.Var().Id("convVal").Op("*").Id(f.table.NameGoLower())
			g.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
			g.Id("c").Dot(f.NameGo()).Op("=").Id("to" + f.table.NameGo() + "Ptr").Call(jen.Id("convVal"))
		})
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("convVal").Id(f.table.NameGoLower())
		g.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convVal"))
		g.Id("c").Dot(f.NameGo()).Op("=").Id("to" + f.table.NameGo()).Call(jen.Id("convVal"))
	})
}
