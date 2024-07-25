package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	filenameInterfaces = "som.interfaces.go"
	filenameSchema     = "tables.surql"
)

type build struct {
	input  *input
	outDir string
	outPkg string
}

func Build(source *parser.Output, outDir string, outPkg string) error {
	in, err := newInput(source)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:  in,
		outDir: outDir,
		outPkg: outPkg,
	}

	return builder.build()
}

func (b *build) build() error {
	if err := os.MkdirAll(b.basePath(), os.ModePerm); err != nil {
		return err
	}

	if err := b.copyInternalPackage(); err != nil {
		return err
	}

	if err := b.embedStaticFiles(); err != nil {
		return err
	}

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

func (b *build) copyInternalPackage() error {
	tmpl := &embed.Template{
		GenerateOutPath: b.subPkg(""),
	}

	files, err := embed.Lib(tmpl)
	if err != nil {
		return err
	}

	dir := filepath.Join(b.outDir, "internal", "lib")

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		err := os.WriteFile(filepath.Join(dir, file.Path), []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *build) embedStaticFiles() error {
	tmpl := &embed.Template{
		GenerateOutPath: b.subPkg(""),
	}

	files, err := embed.Som(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		err := os.WriteFile(filepath.Join(b.basePath(), file.Path), []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *build) buildInterfaceFile() error {
	f := jen.NewFile(b.basePkgName())

	f.PackageComment(codegenComment)

	f.Type().Id("Client").InterfaceFunc(func(g *jen.Group) {
		for _, node := range b.input.nodes {
			g.Id(node.NameGo() + "Repo").Call().Id(node.NameGo() + "Repo")
		}

		g.Id("ApplySchema").Call(jen.Id("ctx").Qual("context", "Context")).Error()
		g.Id("Close").Call()
	})

	if err := f.Save(path.Join(b.basePath(), filenameInterfaces)); err != nil {
		return err
	}

	return nil
}

func (b *build) buildSchemaFile() error {
	statements := []string{"", ""}

	var fieldFn func(table string, f field.Field, prefix string)
	fieldFn = func(table string, f field.Field, prefix string) {
		fieldType := f.TypeDatabase()
		if fieldType == "" {
			return // TODO: is this actually valid?
		}

		statement := fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, fieldType,
		)
		statements = append(statements, statement)

		if object, ok := f.(*field.Struct); ok {
			for _, fld := range object.Table().GetFields() {
				fieldFn(table, fld, prefix+f.NameDatabase()+".")
			}
		}

		if slice, ok := f.(*field.Slice); ok {

			if _, ok := slice.Element().(*field.Byte); ok {
				// byte slice has the type "string" in the database,
				// so we do not need to specify its elements.
				return
			}

			statement := fmt.Sprintf(
				"DEFINE FIELD %s ON TABLE %s TYPE %s;",
				prefix+f.NameDatabase()+".*", table, slice.Element().TypeDatabase(),
			)
			statements = append(statements, statement)

			if object, ok := slice.Element().(*field.Struct); ok {
				for _, fld := range object.Table().GetFields() {
					fieldFn(table, fld, prefix+f.NameDatabase()+".*.")
				}
			}
		}
	}

	for _, node := range b.input.nodes {
		statement := fmt.Sprintf("DEFINE TABLE %s SCHEMAFULL TYPE NORMAL PERMISSIONS FULL;", node.NameDatabase())
		statements = append(statements, statement)

		statement = fmt.Sprintf(
			`DEFINE FIELD id ON TABLE %s TYPE record<%s> ASSERT $value != NONE AND $value != NULL AND $value != "";`,
			node.NameDatabase(), node.NameDatabase(),
		)
		statements = append(statements, statement)

		for _, f := range node.GetFields() {
			fieldFn(node.NameDatabase(), f, "")
		}

		statements = append(statements, "")
	}

	for _, edge := range b.input.edges {
		statement := fmt.Sprintf(
			"DEFINE TABLE %s SCHEMAFULL TYPE RELATION IN %s OUT %s PERMISSIONS FULL;",
			edge.NameDatabase(),
			edge.In.NameDatabase(),
			edge.Out.NameDatabase(), // can be OR'ed with "|"
		)
		statements = append(statements, statement)

		for _, f := range edge.GetFields() {
			fieldFn(edge.NameDatabase(), f, "")
		}

		statements = append(statements, "")
	}

	content := strings.Join(statements, "\n")

	if err := os.MkdirAll(path.Join(b.basePath(), "schema"), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create schema directory: %w", err)
	}

	err := os.WriteFile(path.Join(b.basePath(), "schema", filenameSchema), []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %w", err)
	}

	return nil
}

func (b *build) buildBaseFile(node *field.NodeTable) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(b.basePkgName())

	f.PackageComment(codegenComment)

	//
	// type {NodeName}Repo interface {...}
	//
	f.Type().Id(node.NameGo()+"Repo").Interface(
		jen.Id("Query").Call().Qual(pkgQuery, "Builder").
			Types(b.input.SourceQual(node.NameGo()), jen.Qual(b.subPkg(def.PkgConv), node.NameGo())),

		jen.Id("Create").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("CreateWithID").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Read").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").Op("*").Qual(def.PkgSDBC, "ID"),
		).Parens(jen.List(
			jen.Op("*").Add(b.input.SourceQual(node.NameGo())),
			jen.Bool(),
			jen.Error(),
		)),

		jen.Id("Update").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Delete").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Refresh").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Relate").Call().Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()),
	)

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
								jen.Id("convTo").Op(":").Qual(pkgConv, "To"+node.NameGo()),
							),
							jen.Add(
								jen.Line(),
								jen.Id("convFrom").Op(":").Qual(pkgConv, "From"+node.NameGo()),
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
			jen.Id("id").Op("*").Qual(def.PkgSDBC, "ID"),
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
				),
			),
		)

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

	if err := f.Save(path.Join(b.basePath(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *build) newQueryBuilder() builder {
	return newQueryBuilder(b.input, b.basePath(), b.basePkg(), def.PkgQuery)
}

func (b *build) newFilterBuilder() builder {
	return newFilterBuilder(b.input, b.basePath(), b.basePkg(), def.PkgFilter)
}

func (b *build) newSortBuilder() builder {
	return newSortBuilder(b.input, b.basePath(), b.basePkg(), def.PkgSort)
}

func (b *build) newFetchBuilder() builder {
	return newFetchBuilder(b.input, b.basePath(), b.basePkg(), def.PkgFetch)
}

func (b *build) newConvBuilder() builder {
	return newConvBuilder(b.input, b.basePath(), b.basePkg(), def.PkgConv)
}

func (b *build) newRelateBuilder() builder {
	return newRelateBuilder(b.input, b.basePath(), b.basePkg(), def.PkgRelate)
}

func (b *build) basePath() string {
	return b.outDir
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) basePkgName() string {
	_, name := filepath.Split(b.outPkg)
	return name
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
