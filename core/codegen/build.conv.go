package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type convBuilder struct {
	*baseBuilder
}

func newConvBuilder(input *input, basePath, basePkg, pkgName string) *convBuilder {
	return &convBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *convBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

	// Generate the base file.
	if err := b.buildBaseFile(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildFile(edge); err != nil {
			return err
		}
	}

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *convBuilder) buildBaseFile() error {
	content := `

package conv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func parseDatabaseID(node string, id string) string {
	id = strings.TrimPrefix(id, node+":")
	id = strings.TrimPrefix(id,"⟨")
	id = strings.TrimSuffix(id, "⟩")
	id, _ = strconv.Unquote("\"" + id + "\"")
	return id
}

func buildDatabaseID(node string, id string) string {
	return node + ":" + id
}

func mapEnum[I, O ~string](in I) O {
 	return O(in)
}

func mapSlice[I, O any](in []I, fn func(I) O) []O {
	if in == nil {
		return nil
	}

	out := make([]O, len(in))
	for _, i := range in {
		out = append(out, fn(i))
	}
	return out
}

func mapSlicePtr[I, O any](in *[]I, fn func(I) O) *[]O {
	if in == nil {
		return nil
	}

	out := make([]O, len(*in))
	for _, i := range *in {
		out = append(out, fn(i))
	}
	return &out
}

func mapPtrSlice[I, O any](in []*I, fn func(I) O) []*O {
	if in == nil {
		return nil
	}

 	ptrFn := ptrFunc(fn)

	out := make([]*O, len(in))
 	for _, i := range in {
 		out = append(out, ptrFn(i))
 	}

 	return out
}

func mapPtrSlicePtr[I, O any](in *[]*I, fn func(I) O) *[]*O {
	if in == nil {
		return nil
	}

	ptrFn := ptrFunc(fn)

	out := make([]*O, len(*in))
	for _, i := range *in {
		out = append(out, ptrFn(i))
	}

	return &out
}

func ptrFunc[I, O any](fn func(I) O) func(*I) *O {
 	return func(in *I) *O {
 		if in == nil {
 			return nil
 		}
 		out := fn(*in)
 		return &out
 	}
}

//
// -- UUID
//

type UUID struct {
	*uuid.UUID
}

func (u *UUID) MarshalJSON() ([]byte, error) {
	if u == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(u.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	uid, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("cannot unmarshal uuid: %w", err)
	}

	u.UUID = &uid

	return nil
}
` // TODO: only add uuid functions and import if needed/used

move content above to embed!!!

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "conv.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *convBuilder) buildFile(elem field.Element) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	typeName := elem.NameGoLower()
	if isNode || isEdge {
		typeName = elem.NameGo()
	}

	f.Line()
	f.Type().Id(typeName).StructFunc(func(g *jen.Group) {
		for _, f := range elem.GetFields() {
			if code := f.CodeGen().FieldDef(fieldCtx); code != nil {
				g.Add(code)
			}
		}
	})

	f.Line()
	f.Add(b.buildFrom(elem))

	f.Line()
	f.Add(b.buildTo(elem))

	if node, ok := elem.(*field.NodeTable); ok {
		f.Line()
		f.Type().Id(node.NameGoLower()+"Link").Struct(
			jen.Id(node.NameGo()),
			jen.Id("ID").String(),
		)

		f.Line()
		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGoLower()+"Link")).
			Id("MarshalJSON").Params().
			Params(jen.Index().Byte(), jen.Error()).
			Block(
				jen.If(jen.Id("f").Op("==").Nil()).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				),
				jen.Return(jen.Qual("encoding/json", "Marshal").Call(jen.Id("f").Dot("ID"))),
			)

		f.Line()
		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGoLower()+"Link")).
			Id("UnmarshalJSON").Params(jen.Id("data").Index().Byte()).
			Error().
			Block(
				jen.Id("raw").Op(":=").String().Call(jen.Id("data")),
				jen.If(
					jen.Qual("strings", "HasPrefix").Call(jen.Id("raw"), jen.Lit("\"")).
						Op("&&").Qual("strings", "HasSuffix").Call(jen.Id("raw"), jen.Lit("\"")),
				).
					Block(
						jen.Id("raw").Op("=").Id("raw").Index(jen.Lit(1).Op(":").Len(jen.Id("raw")).Op("-").Lit(1)),
						jen.Id("f").Dot("ID").Op("=").Id("parseDatabaseID").Call(jen.Lit(node.NameDatabase()), jen.Id("raw")),
						jen.Return(jen.Nil()),
					),

				jen.Type().Id("alias").Id(node.NameGoLower()+"Link"),
				jen.Var().Id("link").Id("alias"),

				jen.Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("link")),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.Op("*").Id("f").Op("=").Id(node.NameGoLower()+"Link").Call(jen.Id("link")),
				),

				jen.Return(jen.Err()),
			)

		f.Line()
		f.Add(b.buildFromLink(node))

		f.Line()
		f.Add(b.buildFromLinkPtr(node))

		f.Line()
		f.Add(b.buildToLink(node))

		f.Line()
		f.Add(b.buildToLinkPtr(node))
	}

	if err := f.Save(path.Join(b.path(), elem.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *convBuilder) buildFrom(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	localName := elem.NameGoLower()
	methodPrefix := "from"

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	if isNode || isEdge {
		localName = elem.NameGo()
		methodPrefix = "From"
	}

	return jen.Func().
		Id(methodPrefix + elem.NameGo()).
		Params(jen.Id("data").Add(b.SourceQual(elem.NameGo()))).
		Id(localName).
		Block(
			jen.Return(jen.Id(localName).Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if elem.HasTimestamps() {
						if f.NameGo() == "CreatedAt" || f.NameGo() == "UpdatedAt" {
							continue
						}
					}

					if code := f.CodeGen().ConvFrom(fieldCtx); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}
			}))),
		)
}

func (b *convBuilder) buildTo(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	localName := elem.NameGoLower()
	methodPrefix := "to"

	_, isNode := elem.(*field.NodeTable)
	_, isEdge := elem.(*field.EdgeTable)

	if isNode || isEdge {
		localName = elem.NameGo()
		methodPrefix = "To"
	}

	return jen.Func().
		Id(methodPrefix + elem.NameGo()).
		Params(jen.Id("data").Id(localName)).
		Add(b.SourceQual(elem.NameGo())).
		Block(
			jen.Return(jen.Add(b.SourceQual(elem.NameGo())).Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if elem.HasTimestamps() {
						if f.NameGo() == "CreatedAt" || f.NameGo() == "UpdatedAt" {
							continue
						}
					}

					if code := f.CodeGen().ConvTo(fieldCtx); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}

				if _, ok := elem.(*field.NodeTable); ok {
					d[jen.Id("Node")] = jen.Qual(def.PkgSom, "NewNode").Call(
						jen.Id("parseDatabaseID").Call(
							jen.Lit(elem.NameDatabase()),
							jen.Id("data").Dot("ID"),
						),
					)
				}

				if _, ok := elem.(*field.EdgeTable); ok {
					d[jen.Id("Edge")] = jen.Qual(def.PkgSom, "NewEdge").Call(
						jen.Id("parseDatabaseID").Call(
							jen.Lit(elem.NameDatabase()),
							jen.Id("data").Dot("ID"),
						),
					)
				}

				if elem.HasTimestamps() {
					d[jen.Id("Timestamps")] = jen.Qual(def.PkgSom, "NewTimestamps").Call(
						jen.Id("data").Dot("CreatedAt"),
						jen.Id("data").Dot("UpdatedAt"),
					)
				}
			}))))
}

func (b *convBuilder) buildFromLink(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"Link").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower()+"Link")).
		Add(b.SourceQual(node.NameGo())).
		Block(
			jen.If(jen.Id("link").Op("==").Nil()).Block(
				jen.Return(jen.Add(b.SourceQual(node.NameGo())).Values()),
			),
			jen.Return(jen.Id("To"+node.NameGo()).Call(jen.Id(node.NameGo()).Call(jen.Id("link").Dot(node.NameGo())))),
		)
}

func (b *convBuilder) buildFromLinkPtr(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"LinkPtr").
		Params(jen.Id("link").Op("*").Id(node.NameGoLower()+"Link")).
		Op("*").Add(b.SourceQual(node.NameGo())).
		Block(
			jen.If(jen.Id("link").Op("==").Nil()).Block(
				jen.Return(jen.Nil()),
			),
			jen.Id("node").Op(":=").Id("To"+node.NameGo()).Call(jen.Id(node.NameGo()).Call(jen.Id("link").Dot(node.NameGo()))),
			jen.Return(jen.Op("&").Id("node")),
		)
}

func (b *convBuilder) buildToLink(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("to"+node.NameGo()+"Link").
		Params(jen.Id("node").Add(b.SourceQual(node.NameGo()))).
		Op("*").Id(node.NameGoLower()+"Link").
		Block(
			jen.If(jen.Id("node").Dot("ID").Call().Op("==").Lit("")).Block(
				jen.Return(jen.Nil()),
			),
			jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
				jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(jen.Id("node")),
				jen.Id("ID").Op(":").Id("buildDatabaseID").Call(
					jen.Lit(node.NameDatabase()),
					jen.Id("node").Dot("ID").Call(),
				),
			),
			jen.Return(jen.Op("&").Id("link")),
		)
}

func (b *convBuilder) buildToLinkPtr(node *field.NodeTable) jen.Code {
	return jen.Func().
		Id("to"+node.NameGo()+"LinkPtr").
		Params(jen.Id("node").Op("*").Add(b.SourceQual(node.NameGo()))).
		Op("*").Id(node.NameGoLower()+"Link").
		Block(
			jen.
				If(
					jen.Id("node").Op("==").Nil().Op("||").
						Id("node").Dot("ID").Call().Op("==").Lit(""),
				).
				Block(
					jen.Return(jen.Nil()),
				),
			jen.Id("link").Op(":=").Id(node.NameGoLower()+"Link").Values(
				jen.Id(node.NameGo()).Op(":").Id("From"+node.NameGo()).Call(jen.Op("*").Id("node")),
				jen.Id("ID").Op(":").Id("buildDatabaseID").Call(
					jen.Lit(node.NameDatabase()),
					jen.Id("node").Dot("ID").Call(),
				),
			),
			jen.Return(jen.Op("&").Id("link")),
		)
}
