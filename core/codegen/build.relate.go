package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type relateBuilder struct {
	*baseBuilder
}

func newRelateBuilder(input *input, basePath, basePkg, pkgName string) *relateBuilder {
	return &relateBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *relateBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

	if err := b.buildBaseFile(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildNodeFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildEdgeFile(edge); err != nil {
			return err
		}
	}

	return nil
}

func (b *relateBuilder) buildBaseFile() error {
	content := `package relate

type Database interface {
	Query(statement string, vars map[string]any) (any, error)
}
`

	err := os.WriteFile(path.Join(b.path(), "relate.go"), []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *relateBuilder) buildNodeFile(node *dbtype.Node) error {
	file := jen.NewFile(b.pkgName)

	file.Add(b.byNew(node))

	file.Type().Id(node.Name).Struct(
		jen.Id("db").Id("Database"),
	)

	for _, fld := range node.GetFields() {
		slice, ok := fld.(*field.Slice)
		if !ok {
			continue
		}

		ok, edge, _, _ := slice.Edge()
		if !ok {
			continue
		}

		file.Func().Params(jen.Id("n").Id(node.NameGo())).
			Id(fld.NameGo()).Params().
			Id(strcase.ToLowerCamel(edge)).
			Block(
				jen.Return(jen.Id(strcase.ToLowerCamel(edge)).Values(jen.Id("db").Op(":").Id("n").Dot("db"))),
			)
	}

	if err := file.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) buildEdgeFile(edge *dbtype.Edge) error {
	file := jen.NewFile(b.pkgName)

	file.Type().Id(strcase.ToLowerCamel(edge.Name)).Struct(
		jen.Id("db").Id("Database"),
	)

	in := edge.In.(*field.Node)
	out := edge.Out.(*field.Node)

	file.Func().Params(jen.Id("e").Id(strcase.ToLowerCamel(edge.NameGo()))).
		Id("Create").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.Name))).
		Error().
		Block(
			jen.If(jen.Id("edge").Dot("ID").Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID must not be set for an edge to be created"))),
				),

			jen.If(jen.Id("edge").Dot(in.NameGo()).Dot("ID").Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID of the incoming node '"+in.NameGo()+"' must not be empty"))),
				),

			jen.If(jen.Id("edge").Dot(out.NameGo()).Dot("ID").Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID of the outgoing node '"+out.NameGo()+"' must not be empty"))),
				),

			jen.Id("query").Op(":=").Lit("RELATE "),
			jen.Id("query").Op("+=").Lit(strcase.ToSnake(in.NodeName())+":").Op("+").Id("edge").Dot(in.NameGo()).Dot("ID"),
			jen.Id("query").Op("+=").Lit("->"+edge.NameDatabase()+"->"),
			jen.Id("query").Op("+=").Lit(strcase.ToSnake(out.NodeName())+":").Op("+").Id("edge").Dot(out.NameGo()).Dot("ID"),
			jen.Id("query").Op("+=").Lit(" CONTENT $data"),

			jen.Id("data").Op(":=").Qual(b.subPkg(def.PkgConv), "From"+edge.NameGo()).Call(jen.Id("edge")),
			jen.Id("raw").Op(",").Err().Op(":=").Id("e").Dot("db").Dot("Query").
				Call(jen.Id("query"), jen.Map(jen.String()).Any().Values(jen.Lit("data").Op(":").Id("data"))),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Var().Id("convEdge").Qual(b.subPkg(def.PkgConv), edge.NameGo()),
			jen.List(jen.Id("ok"), jen.Err()).Op(":=").Qual(def.PkgSurrealDB, "UnmarshalRaw").
				Call(jen.Id("raw"), jen.Op("&").Id("convEdge")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("errors", "New").Call(jen.Lit("result is empty"))),
			),

			jen.Op("*").Id("edge").Op("=").Op("*").Qual(b.subPkg(def.PkgConv), "To"+edge.NameGo()).Call(jen.Op("&").Id("convEdge")),
			jen.Return(jen.Nil()),
		)

	//
	//	data := conv.FromMemberOf(edge)
	//	raw, err := e.db.Query(query, map[string]any{"data": data})
	//	if err != nil {
	//		return err
	//	}
	//	var convEdge conv.MemberOf
	//	err = surrealdbgo.Unmarshal(raw, &convEdge)
	//	if err != nil {
	//		return err
	//	}
	//	*edge = *conv.ToMemberOf(&convEdge)
	//	return nil

	file.Func().Params(jen.Id(strcase.ToLowerCamel(edge.NameGo()))).
		Id("Update").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.Name))).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	file.Func().Params(jen.Id(strcase.ToLowerCamel(edge.NameGo()))).
		Id("Delete").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.Name))).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	if err := file.Save(path.Join(b.path(), edge.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) byNew(node dbtype.Element) jen.Code {
	return jen.Func().Id("New" + node.NameGo()).
		Params(jen.Id("db").Id("Database")).
		Id("*").Id(node.NameGo()).
		Block(
			jen.Return(
				jen.Id("&").Id(node.NameGo()).Values(
					jen.Id("db").Op(":").Id("db"),
				),
			),
		)
}
