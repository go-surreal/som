package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

type relateBuilder struct {
	*baseBuilder
}

func newRelateBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *relateBuilder {
	return &relateBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *relateBuilder) build() error {
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

func (b *relateBuilder) buildNodeFile(node *field.NodeTable) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Add(b.byNew(node))

	f.Line()
	f.Type().Id(node.Name).Struct(
		jen.Id("db").Id("Database"),
	)

	for _, fld := range node.GetFields() {
		slice, ok := fld.(*field.Slice)
		if !ok {
			continue
		}

		edgeElement, ok := slice.Element().(*field.Edge)
		if !ok {
			continue
		}

		f.Line()
		f.Func().Params(jen.Id("n").Id(node.NameGo())).
			Id(fld.NameGo()).Params().
			Id(edgeElement.Table().NameGoLower()).
			Block(
				jen.Return(jen.Id(edgeElement.Table().NameGoLower()).Call(jen.Id("n"))),
			)
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) buildEdgeFile(edge *field.EdgeTable) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Type().Id(edge.NameGoLower()).Struct(
		jen.Id("db").Id("Database"),
	)

	f.Line()
	f.Add(
		comment("Create creates a new edge between the given nodes.\nNote: The ID type if both nodes must be a string or number for now."),
	)
	f.Func().Params(jen.Id("e").Id(edge.NameGoLower())).Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("edge").Op("*").Add(b.SourceQual(edge.Name)),
		).
		Error().
		Block(
			jen.If(jen.Id("edge").Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the given edge must not be nil"))),
				),

			jen.If(jen.Id("edge").Dot("ID").Call().Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID must not be set for an edge to be created"))),
				),

			jen.If(jen.Id("edge").Dot(edge.In.NameGo()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID of the incoming node '"+edge.In.NameGo()+"' must not be empty"))),
				),

			jen.If(jen.Id("edge").Dot(edge.Out.NameGo()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID of the outgoing node '"+edge.Out.NameGo()+"' must not be empty"))),
				),

			jen.Id("inID").Op(":=").Qual("github.com/surrealdb/surrealdb.go/pkg/models", "NewRecordID").Call(
				jen.Lit(edge.In.NameDatabase()), b.edgeNodeIDValue(edge.In),
			),
			jen.Id("outID").Op(":=").Qual("github.com/surrealdb/surrealdb.go/pkg/models", "NewRecordID").Call(
				jen.Lit(edge.Out.NameDatabase()), b.edgeNodeIDValue(edge.Out),
			),

			jen.Id("query").Op(":=").Lit("RELATE $inID->"+edge.NameDatabase()+"->$outID CONTENT $data"),

			jen.Id("data").Op(":=").Qual(b.subPkg(def.PkgConv), "From"+edge.NameGo()).Call(jen.Op("*").Id("edge")),

			jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("e").Dot("db").Dot("Query").Call(
				jen.Id("ctx"),
				jen.Id("query"),
				jen.Map(jen.String()).Any().Values(
					jen.Lit("inID").Op(":").Id("inID"),
					jen.Lit("outID").Op(":").Id("outID"),
					jen.Lit("data").Op(":").Id("data"),
				),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not create relation: %w"), jen.Err())),
			),

			jen.Var().Id("rawResult").Index().Qual(b.subPkg(def.PkgInternal), "QueryResult").Types(jen.Qual(b.subPkg(def.PkgConv), edge.NameGo())),
			jen.Err().Op("=").Id("e").Dot("db").Dot("Unmarshal").Call(jen.Id("res"), jen.Op("&").Id("rawResult")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal relation: %w"), jen.Err())),
			),
			jen.If(jen.Len(jen.Id("rawResult")).Op("<").Lit(1).Op("||").Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")).Op("<").Lit(1)).Block(
				jen.Return(jen.Qual("errors", "New").Call(jen.Lit("no result returned for relation"))),
			),

			jen.Id("convEdge").Op(":=").Op("&").Id("rawResult").Index(jen.Lit(0)).Dot("Result").Index(jen.Lit(0)),

			jen.Op("*").Id("edge").Op("=").
				Qual(b.subPkg(def.PkgConv), "To"+edge.NameGo()).Call(jen.Id("convEdge")),

			jen.Return(jen.Nil()),
		)

	f.Line()
	f.Func().Params(jen.Id(edge.NameGoLower())).
		Id("Update").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.NameGo()))).
		Error().
		Block(
			jen.Comment("TODO: implement!"),
			jen.Return(jen.Qual("errors", "New").Call(jen.Lit("not yet implemented"))),
		)

	f.Line()
	f.Func().Params(jen.Id(edge.NameGoLower())).
		Id("Delete").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.NameGo()))).
		Error().
		Block(
			jen.Comment("TODO: implement!"),
			jen.Comment("https://surrealdb.com/docs/surrealdb/surrealql/statements/delete#deleting-graph-edges"),
			jen.Return(jen.Qual("errors", "New").Call(jen.Lit("not yet implemented"))),
		)

	if err := f.Render(b.fs.Writer(path.Join(b.path(), edge.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) edgeNodeIDValue(node *field.Node) jen.Code {
	idExpr := jen.Id("edge").Dot(node.Table().NameGo()).Dot("ID").Call()
	if node.Table().Source.IDType == parser.IDTypeUUID {
		return jen.Qual(b.subPkg(""), "UUID").Call(idExpr)
	}
	return idExpr
}

func (b *relateBuilder) byNew(node field.Element) jen.Code {
	return jen.Func().Id("New" + node.NameGo()).
		Params(
			jen.Id("db").Id("Database"),
		).
		Id("*").Id(node.NameGo()).
		Block(
			jen.Return(
				jen.Id("&").Id(node.NameGo()).Values(
					jen.Id("db").Op(":").Id("db"),
				),
			),
		)
}
