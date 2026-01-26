package codegen

import (
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/embed"
)

func (b *build) buildWireFile() error {
	pkgRepo := b.subPkg(def.PkgRepo)

	f := jen.NewFile(def.PkgSomWire)
	f.PackageComment(string(embed.CodegenComment))

	// var Providers = wire.NewSet(...)
	f.Var().Id("Providers").Op("=").Qual(b.wirePackage, "NewSet").CustomFunc(jen.Options{
		Open:      "(",
		Close:     ")",
		Separator: ",",
		Multi:     true,
	}, func(g *jen.Group) {
		g.Id("ProvideClient")
		g.Qual(b.wirePackage, "Bind").Call(
			jen.New(jen.Qual(pkgRepo, "Client")),
			jen.New(jen.Op("*").Qual(pkgRepo, "ClientImpl")),
		)
		for i, node := range b.input.nodes {
			if i == 0 {
				g.Add(jen.Line(), jen.Id("Provide"+node.NameGo()+"Repo"))
			} else {
				g.Id("Provide" + node.NameGo() + "Repo")
			}
		}
	})

	// func ProvideClient(ctx context.Context, conf repo.Config) (*repo.ClientImpl, func(), error)
	f.Line()
	f.Func().Id("ProvideClient").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("conf").Qual(pkgRepo, "Config"),
	).Params(
		jen.Op("*").Qual(pkgRepo, "ClientImpl"),
		jen.Func().Params(),
		jen.Error(),
	).Block(
		jen.List(jen.Id("client"), jen.Err()).Op(":=").Qual(pkgRepo, "NewClient").Call(jen.Id("ctx"), jen.Id("conf")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Nil(), jen.Err()),
		),
		jen.Id("cleanup").Op(":=").Func().Params().Block(jen.Id("client").Dot("Close").Call()),
		jen.Return(jen.Id("client"), jen.Id("cleanup"), jen.Nil()),
	)

	// Per-node provider functions
	for _, node := range b.input.nodes {
		f.Line()
		f.Func().Id("Provide"+node.NameGo()+"Repo").Params(
			jen.Id("client").Op("*").Qual(pkgRepo, "ClientImpl"),
		).Qual(pkgRepo, node.NameGo()+"Repo").Block(
			jen.Return(jen.Id("client").Dot(node.NameGo() + "Repo").Call()),
		)
	}

	return f.Render(b.fs.Writer(filepath.Join(def.PkgSomWire, "providers.go")))
}
