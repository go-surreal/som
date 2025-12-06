package codegen

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

const filenameInterfaces = "som.interfaces.go"

type build struct {
	input  *input
	fs     *fs.FS
	outPkg string
}

func BuildStatic(fs *fs.FS, outPkg string, features *parser.UsedFeatures) error {
	tmpl := &embed.Template{
		GenerateOutPath: outPkg,
	}

	if features != nil {
		tmpl.UsesGoogleUUID = features.UsesGoogleUUID
		tmpl.UsesGofrsUUID = features.UsesGofrsUUID
	}

	files, err := embed.Read(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		fs.Write(file.Path, file.Content)
	}

	return nil
}

func Build(source *parser.Output, fs *fs.FS, outPkg string) error {
	in, err := newInput(source, outPkg)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:  in,
		fs:     fs,
		outPkg: outPkg,
	}

	return builder.build()
}

func (b *build) build() error {
	if err := b.buildInterfaceFile(); err != nil {
		return err
	}

	if err := b.buildSchemaFile(); err != nil {
		return err
	}

	for _, node := range b.input.nodes {
		if err := b.buildBaseFile(node); err != nil {
			return err
		}
	}

	builders := []builder{
		b.newQueryBuilder(),
		b.newFilterBuilder(),
		b.newSortBuilder(),
		b.newFetchBuilder(),
		b.newConvBuilder(),
		b.newRelateBuilder(),
	}

	for _, builder := range builders {
		if err := builder.build(); err != nil {
			return err
		}
	}

	return nil
}

func (b *build) buildInterfaceFile() error {
	f := jen.NewFile(def.PkgRepo)

	f.PackageComment(string(embed.CodegenComment))

	f.Type().Id("Client").InterfaceFunc(func(g *jen.Group) {
		for _, node := range b.input.nodes {
			g.Id(node.NameGo() + "Repo").Call().Id(node.NameGo() + "Repo")
		}

		g.Id("ApplySchema").Call(jen.Id("ctx").Qual("context", "Context")).Error()
		g.Id("Close").Call()
	})

	if err := f.Render(b.fs.Writer(filepath.Join(def.PkgRepo, filenameInterfaces))); err != nil {
		return err
	}

	return nil
}

func (b *build) buildBaseFile(node *field.NodeTable) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(def.PkgRepo)

	f.PackageComment(string(embed.CodegenComment))

	//
	// type {NodeName}Repo interface {...}
	//
	f.Type().Id(node.NameGo()+"Repo").InterfaceFunc(func(g *jen.Group) {
		g.Id("Query").Call().Qual(pkgQuery, "Builder").
			Types(b.input.SourceQual(node.NameGo()), jen.Qual(b.subPkg(def.PkgConv), node.NameGo()))

		g.Id("Create").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Id("CreateWithID").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Id("Read").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").Op("*").Qual(b.subPkg(""), "ID"),
		).Parens(jen.List(
			jen.Op("*").Add(b.input.SourceQual(node.NameGo())),
			jen.Bool(),
			jen.Error(),
		))

		g.Id("Update").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Id("Delete").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		// Add Purge and Restore for soft delete models
		if node.Source.SoftDelete {
			g.Id("Purge").Call(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
			).Error()

			g.Id("Restore").Call(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
			).Error()
		}

		g.Id("Refresh").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Id("Relate").Call().Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo())
	})

	f.Line().
		Add(comment(`
` + node.NameGo() + `Repo returns a new repository instance for the ` + node.NameGo() + ` model.
		`)).
		Func().Params(jen.Id("c").Op("*").Id("ClientImpl")).
		Id(node.NameGo() + "Repo").Params().Id(node.NameGo() + "Repo").
		Block(
			jen.Return(
				jen.Op("&").Id(node.NameGoLower()).Values(
					jen.Id("repo").Op(":").Op("&").Id("repo").
						Types(
							b.input.SourceQual(node.NameGo()),
							jen.Id("conv."+node.NameGo()),
						).
						Values(
							jen.Add(
								jen.Line(),
								jen.Id("db").Op(":").Id("c").Dot("db"),
							),
							jen.Add(
								jen.Line(),
								jen.Id("name").Op(":").Lit(node.NameDatabase()),
							),
							jen.Add(
								jen.Line(),
								jen.Id("convTo").Op(":").Qual(pkgConv, "To"+node.NameGo()+"Ptr"),
							),
							jen.Add(
								jen.Line(),
								jen.Id("convFrom").Op(":").Qual(pkgConv, "From"+node.NameGo()+"Ptr"),
							),
						),
				),
			),
		)

	f.Line()
	f.Type().Id(node.NameGoLower()).Struct(
		jen.Op("*").Id("repo").Types(
			b.input.SourceQual(node.NameGo()),
			jen.Id("conv."+node.NameGo()),
		),
	)

	f.Line().
		Add(comment(`
Query returns a new query builder for the `+node.NameGo()+` model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Query").Params().
		Qual(pkgQuery, "Builder").
		Types(
			b.input.SourceQual(node.NameGo()),
			jen.Qual(b.subPkg(def.PkgConv), node.NameGo()),
		).
		Block(
			jen.Return(jen.Qual(pkgQuery, "New"+node.NameGo()).Call(
				jen.Id("r").Dot("db"),
			)),
		)

	f.Line().
		Add(comment(`
Create creates a new record for the `+node.NameGo()+` model.
The ID will be generated automatically as a ULID.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))),
				),

			jen.Return(
				jen.Id("r").Dot("create").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
CreateWithID creates a new record for the `+node.NameGo()+` model with the given id.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("CreateWithID").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(), // TODO: name clash if node/model is named "id"!
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))),
				),

			jen.Return(
				jen.Id("r").Dot("createWithID").Call(
					jen.Id("ctx"),
					jen.Id("id"),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Read returns the record for the given id, if it exists.
The returned bool indicates whether the record was found or not.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").Op("*").Qual(b.subPkg(""), "ID"),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Bool(), jen.Error()).
		Block(
			jen.Return(
				jen.Id("r").Dot("read").Call(
					jen.Id("ctx"),
					jen.Id("id"),
				),
			),
		)

	f.Line().
		Add(comment(`
Update updates the record for the given model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Update").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot update "+node.NameGo()+" without existing record ID"))),
				),

			jen.Return(
				jen.Id("r").Dot("update").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Delete deletes the record for the given model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Delete").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				)

			// Check if already deleted (for SoftDelete models)
			if node.Source.SoftDelete {
				g.If(jen.Id(node.NameGoLower()).Dot("SoftDelete").Dot("IsDeleted").Call()).Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("record is already deleted"))),
				)
			}

			g.Return(
				jen.Id("r").Dot("delete").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
					jen.Lit(node.Source.SoftDelete), // Pass softDelete flag
				),
			)
		})

	// Add Purge method for soft delete models
	if node.Source.SoftDelete {
		f.Line().
			Add(comment(`
Purge permanently deletes the record from the database.
This performs a hard delete and cannot be undone.
Use this to permanently remove soft-deleted records.
			`)).
			Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
			Id("Purge").
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
			).
			Error().
			Block(
				jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
					Block(
						jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
					),
				jen.Return(
					jen.Id("r").Dot("delete").Call(
						jen.Id("ctx"),
						jen.Id(node.NameGoLower()).Dot("ID").Call(),
						jen.Id(node.NameGoLower()),
						jen.Lit(false), // Hard delete
					),
				),
			)

		f.Line().
			Add(comment(`
Restore un-deletes a soft-deleted record.
Sets deleted_at to NULL and refreshes the in-memory object.
			`)).
			Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
			Id("Restore").
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
			).
			Error().
			Block(
				jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
					Block(
						jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
					),

				jen.List(jen.Id("patch")).Op(":=").Map(jen.String()).Any().Values(
					jen.Dict{jen.Lit("deleted_at"): jen.Nil()},
				),

				jen.List(jen.Id("_"), jen.Err()).Op(":=").
					Id("r").Dot("db").Dot("Update").Call(
						jen.Id("ctx"),
						jen.Id(node.NameGoLower()).Dot("ID").Call(),
						jen.Id("patch"),
					),

				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not restore entity: %w"), jen.Err())),
				),

				// Auto-refresh to update in-memory object
				jen.Return(jen.Id("r").Dot("refresh").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
				)),
			)
	}

	f.Line().
		Add(comment(`
Refresh refreshes the given model with the remote data.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Refresh").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot refresh "+node.NameGo()+" without existing record ID"))),
				),

			jen.Return(
				jen.Id("r").Dot("refresh").Call(
					jen.Id("ctx"),
					jen.Id(node.NameGoLower()).Dot("ID").Call(),
					jen.Id(node.NameGoLower()),
				),
			),
		)

	f.Line().
		Add(comment(`
Relate returns a new relate instance for the `+node.NameGo()+` model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Relate").Params().
		Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()).
		Block(
			jen.Return(
				jen.Qual(b.subPkg(def.PkgRelate), "New"+node.NameGo()).Call(
					jen.Id("r").Dot("db"),
				),
			),
		)

	if err := f.Render(b.fs.Writer(filepath.Join(def.PkgRepo, node.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *build) newQueryBuilder() builder {
	return newQueryBuilder(b.input, b.fs, b.basePkg(), def.PkgQuery)
}

func (b *build) newFilterBuilder() builder {
	return newFilterBuilder(b.input, b.fs, b.basePkg(), def.PkgFilter)
}

func (b *build) newSortBuilder() builder {
	return newSortBuilder(b.input, b.fs, b.basePkg(), def.PkgSort)
}

func (b *build) newFetchBuilder() builder {
	return newFetchBuilder(b.input, b.fs, b.basePkg(), def.PkgFetch)
}

func (b *build) newConvBuilder() builder {
	return newConvBuilder(b.input, b.fs, b.basePkg(), def.PkgConv)
}

func (b *build) newRelateBuilder() builder {
	return newRelateBuilder(b.input, b.fs, b.basePkg(), def.PkgRelate)
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) subPkg(pkg string) string {
	return path.Join(b.basePkg(), pkg)
}

//
// -- HELPER
//

func comment(text string) jen.Code {
	var code jen.Statement

	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		code.Comment(line).Line()
	}

	return &code
}
