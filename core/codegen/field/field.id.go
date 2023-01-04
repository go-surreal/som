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

func (f *ID) TypeDatabase() string {
	return ""
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
		convTo:   nil, // handled directly, because it is so special ;)
		fieldDef: f.fieldDef,
	}
}

func (f *ID) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLib, "ID").Types(jen.Id("T"))
}

func (f *ID) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLib, "NewID").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Field").Call(jen.Lit(f.NameDatabase())), jen.Lit(ctx.Table.NameDatabase()))
}

func (f *ID) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLib, "BaseSort").Types(jen.Id("T"))
}

func (f *ID) sortInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLib, "NewBaseSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *ID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
