package field

import (
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

func (f *ID) typeConv(_ Context) jen.Code {
	return jen.Op("*").Qual(def.PkgSDBC, "ID") // f.typeGo()
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

func (f *ID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + ",omitempty"})
}
