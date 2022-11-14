package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/parser"
)

type Node struct {
	source          *parser.FieldNode
	dbNameConverter NameConverter
}

// TODO: cool to expose just like that?
func (f *Node) NodeName() string {
	return f.source.Node
}

func (f *Node) NameGo() string {
	return f.source.Name
}

func (f *Node) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Node) FilterDefine(sourcePkg string) jen.Code {
	// Node uses a filter function instead.
	return nil
}

func (f *Node) FilterInit(sourcePkg string, elemName string) jen.Code {
	// Node uses a filter function instead.
	return nil
}

func (f *Node) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(strcase.ToLowerCamel(elem.NameGo())).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strcase.ToLowerCamel(f.source.Node)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Node).Types(jen.Id("T")).
				Params(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Node) SortDefine(types jen.Code) jen.Code {
	// Node uses a sort function instead.
	return nil
}

func (f *Node) SortInit(types jen.Code) jen.Code {
	// Node uses a sort function instead.
	return nil
}

func (f *Node) SortFunc(sourcePkg, elemName string) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(strcase.ToLowerCamel(elemName)).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strcase.ToLowerCamel(f.source.Node)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Node).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Node) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Id("to" + f.source.Node + "Field").Call(jen.Op("&").Id("data").Dot(f.source.Name))
}

func (f *Node) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Op("*").Id("from" + f.source.Node + "Field").Call(jen.Id("data").Dot(f.source.Name))
}

func (f *Node) FieldDef() jen.Code {
	return jen.Id(f.source.Name).Id(f.source.Node + "Field").
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"})
}
