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
	input       *input
	fs          *fs.FS
	outPkg      string
	wirePackage string
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

func Build(source *parser.Output, fs *fs.FS, outPkg string, wirePackage string) error {
	in, err := newInput(source, outPkg)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:       in,
		fs:          fs,
		outPkg:      outPkg,
		wirePackage: wirePackage,
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
		b.newFieldBuilder(),
	}

	for _, builder := range builders {
		if err := builder.build(); err != nil {
			return err
		}
	}

	if b.wirePackage != "" {
		if err := b.buildWireFile(); err != nil {
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

	// Generate ClientImpl with per-node cached repo fields.
	f.Line().Type().Id("ClientImpl").StructFunc(func(g *jen.Group) {
		g.Id("db").Id("Database")
		g.Id("mu").Qual("sync", "Mutex")
		for _, node := range b.input.nodes {
			g.Id(node.NameGoLower() + "Repo").Op("*").Id(node.NameGoLower())
		}
	})

	if err := f.Render(b.fs.Writer(filepath.Join(def.PkgRepo, filenameInterfaces))); err != nil {
		return err
	}

	return nil
}

func (b *build) buildBaseFile(node *field.NodeTable) error {
	pkgQuery := b.relativePkgPath(def.PkgQuery)
	pkgConv := b.relativePkgPath(def.PkgConv)

	f := jen.NewFile(def.PkgRepo)

	f.PackageComment(string(embed.CodegenComment))

	//
	// type {NodeName}Repo interface {...}
	//
	f.Line().Type().Id(node.NameGo()+"Repo").InterfaceFunc(func(g *jen.Group) {
		g.Add(comment("Query returns a new query builder for the " + node.NameGo() + " model."))
		g.Id("Query").Call().Qual(pkgQuery, "Builder").Types(b.input.SourceQual(node.NameGo()))

		g.Add(comment("Create creates a new record for the " + node.NameGo() + " model."))
		g.Id("Create").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Add(comment("CreateWithID creates a new record with the given ID for the " + node.NameGo() + " model."))
		g.Id("CreateWithID").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Add(comment("Read returns the record for the given ID, if it exists."))
		g.Id("Read").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).Parens(jen.List(
			jen.Op("*").Add(b.input.SourceQual(node.NameGo())),
			jen.Bool(),
			jen.Error(),
		))

		g.Add(comment("Update updates the record for the given " + node.NameGo() + " model."))
		g.Id("Update").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Add(comment("Delete deletes the record for the given " + node.NameGo() + " model."))
		g.Id("Delete").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		// Add Erase and Restore for soft delete models
		if node.Source.SoftDelete {
			g.Add(comment("Erase permanently deletes the record from the database."))
			g.Id("Erase").Call(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
			).Error()

			g.Add(comment("Restore un-deletes a soft-deleted record."))
			g.Id("Restore").Call(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
			).Error()
		}

		g.Add(comment("Refresh refreshes the given model with the current database state."))
		g.Id("Refresh").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error()

		g.Add(comment("Relate returns a new relate builder for the " + node.NameGo() + " model."))
		g.Id("Relate").Call().Op("*").Qual(b.relativePkgPath(def.PkgRelate), node.NameGo())

		g.Line()

		for _, event := range []string{"Create", "Update", "Delete"} {
			for _, timing := range []string{"Before", "After"} {
				methodName := "On" + timing + event

				var hookComment string
				switch timing {
				case "Before":
					hookComment = methodName + " registers a hook that runs before a record is " + strings.ToLower(event) + "d.\n" +
						"If the hook returns an error, the " + strings.ToLower(event) + " operation is aborted.\n" +
						"Returns a function that, when called, removes this hook.\n" +
						"\n" +
						"Note: Hooks are local to this application instance and are not\n" +
						"distributed across multiple instances of the application."
				case "After":
					hookComment = methodName + " registers a hook that runs after a record has been " + strings.ToLower(event) + "d.\n" +
						"If the hook returns an error, the error is returned to the caller.\n" +
						"Returns a function that, when called, removes this hook.\n" +
						"\n" +
						"Note: Hooks are local to this application instance and are not\n" +
						"distributed across multiple instances of the application."
				}

				g.Add(comment(hookComment))
				g.Id(methodName).Call(
					jen.Id("fn").Func().Params(
						jen.Id("ctx").Qual("context", "Context"),
						jen.Id("node").Op("*").Add(b.input.SourceQual(node.NameGo())),
					).Error(),
				).Func().Params()
			}
		}
	})

	repoInfoVarName := node.NameGoLower() + "RepoInfo"

	f.Line()
	f.Commentf("%s holds the model-specific conversion functions for %s.", repoInfoVarName, node.NameGo())
	f.Var().Id(repoInfoVarName).Op("=").Id("RepoInfo").Types(b.input.SourceQual(node.NameGo())).Values(jen.Dict{
		jen.Id("UnmarshalOne"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
		).Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Error()).Block(
			jen.Var().Id("raw").Op("*").Qual(pkgConv, node.NameGo()),
			jen.If(jen.Err().Op(":=").Id("unmarshal").Call(jen.Id("data"), jen.Op("&").Id("raw")), jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Qual(pkgConv, "To"+node.NameGo()+"Ptr").Call(jen.Id("raw")), jen.Nil()),
		),
		jen.Id("MarshalOne"): jen.Func().Params(
			jen.Id("node").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Any().Block(
			jen.Return(jen.Qual(pkgConv, "From"+node.NameGo()+"Ptr").Call(jen.Id("node"))),
		),
	})

	f.Line().
		Add(comment(`
` + node.NameGo() + `Repo returns the repository instance for the ` + node.NameGo() + ` model.
The instance is cached as a singleton on the client.
		`)).
		Func().Params(jen.Id("c").Op("*").Id("ClientImpl")).
		Id(node.NameGo() + "Repo").Params().Id(node.NameGo() + "Repo").
		Block(
			jen.Id("c").Dot("mu").Dot("Lock").Call(),
			jen.Defer().Id("c").Dot("mu").Dot("Unlock").Call(),
			jen.If(jen.Id("c").Dot(node.NameGoLower()+"Repo").Op("==").Nil()).Block(
				jen.Id("c").Dot(node.NameGoLower()+"Repo").Op("=").
					Op("&").Id(node.NameGoLower()).Values(
					jen.Id("repo").Op(":").Op("&").Id("repo").
						Types(
							b.input.SourceQual(node.NameGo()),
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
								jen.Id("info").Op(":").Id(repoInfoVarName),
							),
							jen.Add(
								jen.Line(),
								jen.Id("newID").Op(":").Id(idFuncName(node)),
							),
							jen.Add(
								jen.Line(),
								jen.Id("parseID").Op(":").Id(parseIDFuncName(node)),
							),
						),
				),
			),
			jen.Return(jen.Id("c").Dot(node.NameGoLower()+"Repo")),
		)

	hookEvents := []string{"Create", "Update", "Delete"}

	f.Line()
	f.Type().Id(node.NameGoLower()).StructFunc(func(g *jen.Group) {
		g.Op("*").Id("repo").Types(b.input.SourceQual(node.NameGo()))
		g.Id("mu").Qual("sync", "RWMutex")
		for _, event := range hookEvents {
			for _, timing := range []string{"before", "after"} {
				g.Id(timing + event).Index().Id(node.NameGoLower() + "Hook")
			}
		}
	})

	f.Line()
	f.Type().Id(node.NameGoLower() + "Hook").Struct(
		jen.Id("id").Uint64(),
		jen.Id("fn").Func().Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("node").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),
	)

	f.Line()
	f.Var().Id(node.NameGoLower() + "HookCounter").Qual("sync/atomic", "Uint64")

	for _, event := range hookEvents {
		for _, timing := range []string{"Before", "After"} {
			methodName := "On" + timing + event
			fieldName := strings.ToLower(timing) + event

			var hookComment string
			switch timing {
			case "Before":
				hookComment = methodName + " registers a hook that runs before a record is " + strings.ToLower(event) + "d.\n" +
					"If the hook returns an error, the " + strings.ToLower(event) + " operation is aborted.\n" +
					"Returns a function that, when called, removes this hook.\n" +
					"\n" +
					"Note: Hooks are local to this application instance and are not\n" +
					"distributed across multiple instances of the application."
			case "After":
				hookComment = methodName + " registers a hook that runs after a record has been " + strings.ToLower(event) + "d.\n" +
					"If the hook returns an error, the error is returned to the caller.\n" +
					"Returns a function that, when called, removes this hook.\n" +
					"\n" +
					"Note: Hooks are local to this application instance and are not\n" +
					"distributed across multiple instances of the application."
			}

			f.Line().
				Add(comment(hookComment)).
				Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
				Id(methodName).
				Params(jen.Id("fn").Func().Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("node").Op("*").Add(b.input.SourceQual(node.NameGo())),
				).Error()).
				Func().Params().
				Block(
					jen.Id("id").Op(":=").Id(node.NameGoLower()+"HookCounter").Dot("Add").Call(jen.Lit(1)),
					jen.Id("r").Dot("mu").Dot("Lock").Call(),
					jen.Id("r").Dot(fieldName).Op("=").Append(
						jen.Id("r").Dot(fieldName),
						jen.Id(node.NameGoLower()+"Hook").Values(jen.Dict{
							jen.Id("id"): jen.Id("id"),
							jen.Id("fn"): jen.Id("fn"),
						}),
					),
					jen.Id("r").Dot("mu").Dot("Unlock").Call(),
					jen.Return(jen.Func().Params().Block(
						jen.Id("r").Dot("mu").Dot("Lock").Call(),
						jen.Defer().Id("r").Dot("mu").Dot("Unlock").Call(),
						jen.For(jen.Id("i").Op(",").Id("h").Op(":=").Range().Id("r").Dot(fieldName)).Block(
							jen.If(jen.Id("h").Dot("id").Op("==").Id("id")).Block(
								jen.Id("r").Dot(fieldName).Op("=").Qual("slices", "Delete").Call(
									jen.Id("r").Dot(fieldName),
									jen.Id("i"),
									jen.Id("i").Op("+").Lit(1),
								),
								jen.Return(),
							),
						),
					)),
				)
		}
	}

	f.Line().
		Add(comment(`
Query returns a new query builder for the `+node.NameGo()+` model.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Query").Params().
		Qual(pkgQuery, "Builder").Types(b.input.SourceQual(node.NameGo())).
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
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))))
			g.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Lit("")).
				Block(jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))))

			b.addBeforeHooks(g, node, "Create")

			g.If(jen.Err().Op(":=").Id("r").Dot("create").Call(
				jen.Id("ctx"), jen.Id(node.NameGoLower()),
			), jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))

			b.addAfterHooks(g, node, "Create")

			g.Return(jen.Nil())
		})

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
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))))
			g.If(jen.Id("id").Op("==").Lit("")).
				Block(jen.Return(jen.Qual(b.relativePkgPath(), "ErrEmptyID")))
			g.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Lit("")).
				Block(jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))))

			b.addBeforeHooks(g, node, "Create")

			g.If(jen.Err().Op(":=").Id("r").Dot("createWithID").Call(
				jen.Id("ctx"), jen.Id("id"), jen.Id(node.NameGoLower()),
			), jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))

			b.addAfterHooks(g, node, "Create")

			g.Return(jen.Nil())
		})

	f.Line().
		Add(comment(`
Read returns the record for the given id, if it exists.
The returned bool indicates whether the record was found or not.
If caching is enabled via som.WithCache, it will be used.
		`)).
		Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Bool(), jen.Error()).
		Block(
			jen.If(jen.Id("id").Op("==").Lit("")).Block(
				jen.Return(jen.Nil(), jen.False(), jen.Qual(b.relativePkgPath(), "ErrEmptyID")),
			),
			jen.If(jen.Op("!").Qual(b.relativePkgPath("internal"), "CacheEnabled").Types(b.input.SourceQual(node.NameGo())).Call(jen.Id("ctx"))).Block(
				jen.Return(jen.Id("r").Dot("read").Call(jen.Id("ctx"), jen.Id("r").Dot("recordID").Call(jen.Id("id")))),
			),
			jen.Id("idFunc").Op(":=").Func().Params(jen.Id("n").Op("*").Add(b.input.SourceQual(node.NameGo()))).String().Block(
				jen.Return(jen.Id("n").Dot("ID").Call()),
			),
			jen.Id("queryAll").Op(":=").Func().Params(jen.Id("ctx").Qual("context", "Context")).Params(jen.Index().Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Error()).Block(
				jen.Return(jen.Id("r").Dot("Query").Call().Dot("All").Call(jen.Id("ctx"))),
			),
			jen.Id("countAll").Op(":=").Func().Params(jen.Id("ctx").Qual("context", "Context")).Params(jen.Int(), jen.Error()).Block(
				jen.Return(jen.Id("r").Dot("Query").Call().Dot("Count").Call(jen.Id("ctx"))),
			),
			jen.List(jen.Id("cache"), jen.Err()).Op(":=").Id("getOrCreateCache").
				Types(b.input.SourceQual(node.NameGo())).
				Call(
					jen.Id("ctx"),
					jen.Id("idFunc"),
					jen.Id("queryAll"),
					jen.Id("countAll"),
				),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.False(), jen.Err()),
			),
			jen.Var().Id("refreshFuncs").Op("*").Id("eagerRefreshFuncs").Types(b.input.SourceQual(node.NameGo())),
			jen.If(jen.Id("cache").Op("!=").Nil().Op("&&").Id("cache").Dot("isEager").Call()).Block(
				jen.Id("refreshFuncs").Op("=").Op("&").Id("eagerRefreshFuncs").Types(b.input.SourceQual(node.NameGo())).Values(
					jen.Id("cacheID").Op(":").Qual(b.relativePkgPath("internal"), "GetCacheKey").Types(b.input.SourceQual(node.NameGo())).Call(jen.Id("ctx")),
					jen.Id("queryAll").Op(":").Id("queryAll"),
					jen.Id("countAll").Op(":").Id("countAll"),
					jen.Id("idFunc").Op(":").Id("idFunc"),
				),
			),
			jen.Return(
				jen.Id("r").Dot("readWithCache").Call(
					jen.Id("ctx"),
					jen.Id("id"),
					jen.Id("cache"),
					jen.Id("refreshFuncs"),
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
		BlockFunc(func(g *jen.Group) {
			g.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))))
			g.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Lit("")).
				Block(jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot update "+node.NameGo()+" without existing record ID"))))

			b.addBeforeHooks(g, node, "Update")

			g.If(jen.Err().Op(":=").Id("r").Dot("update").Call(
				jen.Id("ctx"), jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()), jen.Id(node.NameGoLower()),
			), jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))

			b.addAfterHooks(g, node, "Update")

			g.Return(jen.Nil())
		})

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

			g.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot delete "+node.NameGo()+" without existing record ID"))),
				)

			if node.Source.SoftDelete {
				g.If(jen.Id(node.NameGoLower()).Dot("SoftDelete").Dot("IsDeleted").Call()).Block(
					jen.Return(jen.Qual(b.relativePkgPath(), "ErrAlreadyDeleted")),
				)
			}

			b.addBeforeHooks(g, node, "Delete")

			if node.Source.SoftDelete && node.Source.OptimisticLock {
				g.Id("version").Op(":=").Id(node.NameGoLower()).Dot("Version").Call()
				g.If(jen.Err().Op(":=").Id("r").Dot("delete").Call(
					jen.Id("ctx"),
					jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()),
					jen.Id(node.NameGoLower()),
					jen.Lit(true),
					jen.Op("&").Id("version"),
				), jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
			} else {
				g.If(jen.Err().Op(":=").Id("r").Dot("delete").Call(
					jen.Id("ctx"),
					jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()),
					jen.Id(node.NameGoLower()),
					jen.Lit(node.Source.SoftDelete),
					jen.Nil(),
				), jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
			}

			b.addAfterHooks(g, node, "Delete")

			g.Return(jen.Nil())
		})

	// Add Erase method for soft delete models
	if node.Source.SoftDelete {
		f.Line().
			Add(comment(`
Erase permanently deletes the record from the database.
This performs a hard delete and cannot be undone.
Use this to permanently remove soft-deleted records.
			`)).
			Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
			Id("Erase").
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
				jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Lit("")).
					Block(
						jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot erase "+node.NameGo()+" without existing record ID"))),
					),
				jen.Return(
					jen.Id("r").Dot("delete").Call(
						jen.Id("ctx"),
						jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()),
						jen.Id(node.NameGoLower()),
						jen.Lit(false), // Hard delete
						jen.Nil(),
					),
				),
			)

		f.Line().
			Add(comment(`
Restore un-deletes a soft-deleted record.
Sets deleted_at to NONE and refreshes the in-memory object.
			`)).
			Func().Params(jen.Id("r").Op("*").Id(node.NameGoLower())).
			Id("Restore").
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

				g.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Lit("")).Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot restore "+node.NameGo()+" without existing record ID"))),
				)

				g.If(jen.Op("!").Id(node.NameGoLower()).Dot("SoftDelete").Dot("IsDeleted").Call()).Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("record is not deleted, cannot restore"))),
				)

				if node.Source.OptimisticLock {
					g.Add(jen.Id("query").Op(":=").Lit("UPDATE $id SET deleted_at = NONE, __som_lock_version = $lock_version"))
					g.Add(jen.Id("vars").Op(":=").Map(jen.String()).Any().Values(
						jen.Dict{
							jen.Lit("id"):           jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()),
							jen.Lit("lock_version"): jen.Id(node.NameGoLower()).Dot("Version").Call(),
						},
					))
				} else {
					g.Add(jen.Id("query").Op(":=").Lit("UPDATE $id SET deleted_at = NONE"))
					g.Add(jen.Id("vars").Op(":=").Map(jen.String()).Any().Values(
						jen.Dict{jen.Lit("id"): jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call())},
					))
				}

				g.List(jen.Id("_"), jen.Err()).Op(":=").
					Id("r").Dot("db").Dot("Query").Call(
						jen.Id("ctx"),
						jen.Id("query"),
						jen.Id("vars"),
					)

				if node.Source.OptimisticLock {
					g.If(jen.Err().Op("!=").Nil()).Block(
						jen.If(jen.Qual("strings", "Contains").Call(
							jen.Err().Dot("Error").Call(), jen.Lit("optimistic_lock_failed"),
						)).Block(
							jen.Return(jen.Qual(b.relativePkgPath(), "ErrOptimisticLock")),
						),
						jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not restore entity: %w"), jen.Err())),
					)
				} else {
					g.If(jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not restore entity: %w"), jen.Err())),
					)
				}

				g.Return(jen.Id("r").Dot("refresh").Call(
					jen.Id("ctx"),
					jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()),
					jen.Id(node.NameGoLower()),
				))
			})
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

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot refresh "+node.NameGo()+" without existing record ID"))),
				),

			jen.Return(
				jen.Id("r").Dot("refresh").Call(
					jen.Id("ctx"),
					jen.Id("r").Dot("recordID").Call(jen.Id(node.NameGoLower()).Dot("ID").Call()),
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
		Op("*").Qual(b.relativePkgPath(def.PkgRelate), node.NameGo()).
		Block(
			jen.Return(
				jen.Qual(b.relativePkgPath(def.PkgRelate), "New"+node.NameGo()).Call(
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

func (b *build) newFieldBuilder() builder {
	return newFieldBuilder(b.input, b.fs, b.basePkg(), def.PkgField)
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) relativePkgPath(pkg ...string) string {
	return path.Join(append([]string{b.basePkg()}, pkg...)...)
}

//
// -- HELPER
//

func (b *build) addBeforeHooks(g *jen.Group, node *field.NodeTable, event string) {
	somPkg := b.relativePkgPath()
	hookIface := "OnBefore" + event + "Hook"
	fieldName := "before" + event

	g.If(
		jen.List(jen.Id("h"), jen.Id("ok")).Op(":=").Any().Call(jen.Id(node.NameGoLower())).Assert(jen.Qual(somPkg, hookIface)),
		jen.Id("ok"),
	).Block(
		jen.If(jen.Err().Op(":=").Id("h").Dot("OnBefore"+event).Call(jen.Id("ctx")), jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
	)

	g.Id("r").Dot("mu").Dot("RLock").Call()
	g.Id(fieldName+"Hooks").Op(":=").Make(jen.Index().Id(node.NameGoLower()+"Hook"), jen.Len(jen.Id("r").Dot(fieldName)))
	g.Copy(jen.Id(fieldName+"Hooks"), jen.Id("r").Dot(fieldName))
	g.Id("r").Dot("mu").Dot("RUnlock").Call()
	g.For(jen.List(jen.Id("_"), jen.Id("h")).Op(":=").Range().Id(fieldName + "Hooks")).Block(
		jen.If(jen.Err().Op(":=").Id("h").Dot("fn").Call(jen.Id("ctx"), jen.Id(node.NameGoLower())), jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
	)
}

func (b *build) addAfterHooks(g *jen.Group, node *field.NodeTable, event string) {
	somPkg := b.relativePkgPath()
	hookIface := "OnAfter" + event + "Hook"
	fieldName := "after" + event

	g.If(
		jen.List(jen.Id("h"), jen.Id("ok")).Op(":=").Any().Call(jen.Id(node.NameGoLower())).Assert(jen.Qual(somPkg, hookIface)),
		jen.Id("ok"),
	).Block(
		jen.If(jen.Err().Op(":=").Id("h").Dot("OnAfter"+event).Call(jen.Id("ctx")), jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
	)

	g.Id("r").Dot("mu").Dot("RLock").Call()
	g.Id(fieldName+"Hooks").Op(":=").Make(jen.Index().Id(node.NameGoLower()+"Hook"), jen.Len(jen.Id("r").Dot(fieldName)))
	g.Copy(jen.Id(fieldName+"Hooks"), jen.Id("r").Dot(fieldName))
	g.Id("r").Dot("mu").Dot("RUnlock").Call()
	g.For(jen.List(jen.Id("_"), jen.Id("h")).Op(":=").Range().Id(fieldName + "Hooks")).Block(
		jen.If(jen.Err().Op(":=").Id("h").Dot("fn").Call(jen.Id("ctx"), jen.Id(node.NameGoLower())), jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
	)
}

func parseIDFuncName(node *field.NodeTable) string {
	switch node.Source.IDType {
	case parser.IDTypeUUID:
		return "parseUUID"
	default:
		return "parseStringID"
	}
}

func idFuncName(node *field.NodeTable) string {
	switch node.Source.IDType {
	case parser.IDTypeUUID:
		return "newUUID"
	case parser.IDTypeRand:
		return "newID"
	default:
		return "newULID" // ULID is the default ID type (used by the Node alias)
	}
}

func comment(text string) jen.Code {
	var code jen.Statement

	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		code.Comment(line).Line()
	}

	return &code
}
