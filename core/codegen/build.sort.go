package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
	"path/filepath"
	"strings"
)

type sortBuilder struct {
	*baseBuilder
}

func newSortBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *sortBuilder {
	return &sortBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *sortBuilder) build() error {
	if err := b.embedStaticFiles(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *sortBuilder) embedStaticFiles() error {
	tmpl := &embed.Template{
		GenerateOutPath: b.subPkg(""),
	}

	files, err := embed.Sort(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		b.fs.Write(filepath.Join(b.path(), file.Path), []byte(content))
	}

	return nil
}

func (b *sortBuilder) buildFile(node *field.NodeTable) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     node,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	f.Line()
	f.Var().Id(node.Name).Op("=").Id("new" + node.Name).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Line()
	f.Add(b.byNew(node))

	f.Line()
	f.Type().Id(node.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").String())
			for _, f := range node.Fields {
				if code := f.CodeGen().SortDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range node.GetFields() {
		if code := fld.CodeGen().SortFunc(fieldCtx); code != nil {
			f.Line()
			f.Add(code)
		}
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *sortBuilder) byNew(node *field.NodeTable) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     node,
	}

	return jen.Func().Id("new" + node.Name).
		Types(jen.Add(def.TypeModel).Any()).
		Params(jen.Id("key").String()).
		Id(node.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(
				jen.Id(node.NameGoLower()).Types(def.TypeModel).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("key")] = jen.Id("key")
						for _, f := range node.Fields {
							if code := f.CodeGen().SortInit(fieldCtx); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}
