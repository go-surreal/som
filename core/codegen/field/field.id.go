package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type ID struct {
	*baseField

	source *parser.FieldID
}

func (f *ID) typeGo() jen.Code {
	return jen.String()
}

func (f *ID) typeConv() jen.Code {
	return f.typeGo()
}

func (f *ID) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		convFrom: nil, // the ID field must not be passed as field
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *ID) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "ID").Types(jen.Id("T"))
}

func (f *ID) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewID").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(f.NameDatabase())), jen.Lit(ctx.Table.NameDatabase()))
}

func (f *ID) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibSort, "Sort").Types(jen.Id("T"))
}

func (f *ID) sortInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *ID) convTo(ctx Context) jen.Code {
	return jen.Id(funcParseDatabaseID).Call(
		jen.Lit(ctx.Table.NameDatabase()),
		jen.Id("data").Dot(f.NameGo()),
	)
}

func (f *ID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}