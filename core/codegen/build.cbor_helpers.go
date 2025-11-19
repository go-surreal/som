package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
)

type cborHelpersBuilder struct {
	*baseBuilder
}

func newCBORHelpersBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *cborHelpersBuilder {
	return &cborHelpersBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, def.PkgCBORHelpers),
	}
}

func (b *cborHelpersBuilder) build() error {
	f := jen.NewFile("cbor")
	f.PackageComment(string(embed.CodegenComment))

	f.ImportAlias(def.PkgCBOR, "cbor")
	f.ImportName("time", "time")
	f.ImportName("net/url", "url")
	f.ImportName("github.com/google/uuid", "uuid")
	f.ImportName("github.com/surrealdb/surrealdb.go/pkg/models", "models")

	// Constants
	f.Line()
	f.Const().Defs(
		jen.Id("tagDatetime").Op("=").Lit(12),
		jen.Id("tagDuration").Op("=").Lit(14),
		jen.Id("nanosecond").Op("=").Lit(1e9),
	)

	// DateTime helpers
	f.Line()
	f.Comment("DateTime marshaling helpers")
	f.Add(b.buildMarshalDateTime())
	f.Line()
	f.Add(b.buildMarshalDateTimePtr())
	f.Line()
	f.Add(b.buildUnMarshalDateTime())
	f.Line()
	f.Add(b.buildUnMarshalDateTimePtr())

	// Duration helpers
	f.Line()
	f.Comment("Duration marshaling helpers")
	f.Add(b.buildMarshalDuration())
	f.Line()
	f.Add(b.buildMarshalDurationPtr())
	f.Line()
	f.Add(b.buildUnMarshalDuration())
	f.Line()
	f.Add(b.buildUnMarshalDurationPtr())

	// UUID helpers
	f.Line()
	f.Comment("UUID marshaling helpers")
	f.Add(b.buildMarshalUUID())
	f.Line()
	f.Add(b.buildMarshalUUIDPtr())
	f.Line()
	f.Add(b.buildUnMarshalUUID())
	f.Line()
	f.Add(b.buildUnMarshalUUIDPtr())

	// URL helpers
	f.Line()
	f.Comment("URL marshaling helpers")
	f.Add(b.buildMarshalURL())
	f.Line()
	f.Add(b.buildMarshalURLPtr())
	f.Line()
	f.Add(b.buildUnMarshalURL())
	f.Line()
	f.Add(b.buildUnMarshalURLPtr())

	return f.Render(b.fs.Writer(path.Join(b.path(), "helpers.go")))
}

// DateTime helpers

func (b *cborHelpersBuilder) buildMarshalDateTime() jen.Code {
	return jen.Func().Id("MarshalDateTime").
		Params(jen.Id("t").Qual("time", "Time")).
		Params(jen.Qual(def.PkgCBOR, "RawMessage"), jen.Error()).
		Block(
			jen.Id("content").Op(",").Id("err").Op(":=").Qual(def.PkgCBOR, "Marshal").Call(
				jen.Index().Int64().Values(
					jen.Id("t").Dot("Unix").Call(),
					jen.Int64().Call(jen.Id("t").Dot("Nanosecond").Call()),
				),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(
				jen.Qual(def.PkgCBOR, "Marshal").Call(
					jen.Qual(def.PkgCBOR, "RawTag").Values(jen.Dict{
						jen.Id("Number"):  jen.Id("tagDatetime"),
						jen.Id("Content"): jen.Id("content"),
					}),
				),
			),
		)
}

func (b *cborHelpersBuilder) buildMarshalDateTimePtr() jen.Code {
	return jen.Func().Id("MarshalDateTimePtr").
		Params(jen.Id("t").Op("*").Qual("time", "Time")).
		Params(jen.Qual(def.PkgCBOR, "RawMessage"), jen.Error()).
		Block(
			jen.If(jen.Id("t").Op("==").Nil()).Block(
				jen.Return(jen.Qual(def.PkgCBOR, "Marshal").Call(jen.Nil())),
			),
			jen.Return(jen.Id("MarshalDateTime").Call(jen.Op("*").Id("t"))),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalDateTime() jen.Code {
	return jen.Func().Id("UnmarshalDateTime").
		Params(jen.Id("data").Index().Byte()).
		Params(jen.Qual("time", "Time"), jen.Error()).
		Block(
			jen.Var().Id("val").Index().Int64(),
			jen.If(
				jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(
					jen.Id("data"),
					jen.Op("&").Id("val"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Qual("time", "Time").Values(), jen.Err()),
			),
			jen.If(jen.Len(jen.Id("val")).Op("==").Lit(0)).Block(
				jen.Return(jen.Qual("time", "Time").Values(), jen.Nil()),
			),
			jen.Id("secs").Op(":=").Id("val").Index(jen.Lit(0)),
			jen.Id("nano").Op(":=").Int64().Call(jen.Lit(0)),
			jen.If(jen.Len(jen.Id("val")).Op(">").Lit(1)).Block(
				jen.Id("nano").Op("=").Id("val").Index(jen.Lit(1)),
			),
			jen.Return(jen.Qual("time", "Unix").Call(jen.Id("secs"), jen.Id("nano")), jen.Nil()),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalDateTimePtr() jen.Code {
	return jen.Func().Id("UnmarshalDateTimePtr").
		Params(jen.Id("data").Index().Byte()).
		Params(jen.Op("*").Qual("time", "Time"), jen.Error()).
		Block(
			jen.Id("t").Op(",").Id("err").Op(":=").Id("UnmarshalDateTime").Call(jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Op("&").Id("t"), jen.Nil()),
		)
}

// Duration helpers

func (b *cborHelpersBuilder) buildMarshalDuration() jen.Code {
	return jen.Func().Id("MarshalDuration").
		Params(jen.Id("d").Qual("time", "Duration")).
		Params(jen.Qual(def.PkgCBOR, "RawMessage"), jen.Error()).
		Block(
			jen.Id("totalSeconds").Op(":=").Int64().Call(jen.Id("d").Dot("Seconds").Call()),
			jen.Id("totalNanoseconds").Op(":=").Id("d").Dot("Nanoseconds").Call(),
			jen.Id("remainingNanoseconds").Op(":=").Id("totalNanoseconds").Op("-").Parens(jen.Id("totalSeconds").Op("*").Id("nanosecond")),
			jen.Id("content").Op(",").Id("err").Op(":=").Qual(def.PkgCBOR, "Marshal").Call(
				jen.Index().Int64().Values(jen.Id("totalSeconds"), jen.Id("remainingNanoseconds")),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(
				jen.Qual(def.PkgCBOR, "Marshal").Call(
					jen.Qual(def.PkgCBOR, "RawTag").Values(jen.Dict{
						jen.Id("Number"):  jen.Id("tagDuration"),
						jen.Id("Content"): jen.Id("content"),
					}),
				),
			),
		)
}

func (b *cborHelpersBuilder) buildMarshalDurationPtr() jen.Code {
	return jen.Func().Id("MarshalDurationPtr").
		Params(jen.Id("d").Op("*").Qual("time", "Duration")).
		Params(jen.Qual(def.PkgCBOR, "RawMessage"), jen.Error()).
		Block(
			jen.If(jen.Id("d").Op("==").Nil()).Block(
				jen.Return(jen.Qual(def.PkgCBOR, "Marshal").Call(jen.Nil())),
			),
			jen.Return(jen.Id("MarshalDuration").Call(jen.Op("*").Id("d"))),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalDuration() jen.Code {
	return jen.Func().Id("UnmarshalDuration").
		Params(jen.Id("data").Index().Byte()).
		Params(jen.Qual("time", "Duration"), jen.Error()).
		Block(
			jen.Var().Id("val").Index().Int64(),
			jen.If(
				jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(
					jen.Id("data"),
					jen.Op("&").Id("val"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Lit(0), jen.Err()),
			),
			jen.Var().Id("dur").Qual("time", "Duration"),
			jen.If(jen.Len(jen.Id("val")).Op(">").Lit(0)).Block(
				jen.Id("dur").Op("=").Qual("time", "Duration").Call(jen.Id("val").Index(jen.Lit(0))).Op("*").Qual("time", "Second"),
			),
			jen.If(jen.Len(jen.Id("val")).Op(">").Lit(1)).Block(
				jen.Id("dur").Op("+=").Qual("time", "Duration").Call(jen.Id("val").Index(jen.Lit(1))),
			),
			jen.Return(jen.Id("dur"), jen.Nil()),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalDurationPtr() jen.Code {
	return jen.Func().Id("UnmarshalDurationPtr").
		Params(jen.Id("data").Index().Byte()).
		Params(jen.Op("*").Qual("time", "Duration"), jen.Error()).
		Block(
			jen.Id("d").Op(",").Id("err").Op(":=").Id("UnmarshalDuration").Call(jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Op("&").Id("d"), jen.Nil()),
		)
}

// UUID helpers

func (b *cborHelpersBuilder) buildMarshalUUID() jen.Code {
	return jen.Func().Id("MarshalUUID").
		Params(jen.Id("u").Qual("github.com/google/uuid", "UUID")).
		Params(jen.Qual(def.PkgCBOR, "RawMessage"), jen.Error()).
		Block(
			jen.Id("raw").Op(",").Id("err").Op(":=").Qual(def.PkgCBOR, "Marshal").Call(jen.Id("u")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(
				jen.Qual(def.PkgCBOR, "Marshal").Call(
					jen.Qual(def.PkgCBOR, "RawTag").Values(jen.Dict{
						jen.Id("Number"):  jen.Qual("github.com/surrealdb/surrealdb.go/pkg/models", "TagSpecBinaryUUID"),
						jen.Id("Content"): jen.Id("raw"),
					}),
				),
			),
		)
}

func (b *cborHelpersBuilder) buildMarshalUUIDPtr() jen.Code {
	return jen.Func().Id("MarshalUUIDPtr").
		Params(jen.Id("u").Op("*").Qual("github.com/google/uuid", "UUID")).
		Params(jen.Qual(def.PkgCBOR, "RawMessage"), jen.Error()).
		Block(
			jen.If(jen.Id("u").Op("==").Nil()).Block(
				jen.Return(jen.Qual(def.PkgCBOR, "Marshal").Call(jen.Nil())),
			),
			jen.Return(jen.Id("MarshalUUID").Call(jen.Op("*").Id("u"))),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalUUID() jen.Code {
	return jen.Func().Id("UnmarshalUUID").
		Params(jen.Id("data").Index().Byte()).
		Params(jen.Qual("github.com/google/uuid", "UUID"), jen.Error()).
		Block(
			jen.Var().Id("tag").Qual(def.PkgCBOR, "RawTag"),
			jen.If(
				jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(
					jen.Id("data"),
					jen.Op("&").Id("tag"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Qual("github.com/google/uuid", "UUID").Values(), jen.Err()),
			),
			jen.Var().Id("uuidBytes").Index().Byte(),
			jen.If(
				jen.Err().Op(":=").Qual(def.PkgCBOR, "Unmarshal").Call(
					jen.Id("tag").Dot("Content"),
					jen.Op("&").Id("uuidBytes"),
				),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Return(jen.Qual("github.com/google/uuid", "UUID").Values(), jen.Err()),
			),
			jen.Return(jen.Qual("github.com/google/uuid", "FromBytes").Call(jen.Id("uuidBytes"))),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalUUIDPtr() jen.Code {
	return jen.Func().Id("UnmarshalUUIDPtr").
		Params(jen.Id("data").Index().Byte()).
		Params(jen.Op("*").Qual("github.com/google/uuid", "UUID"), jen.Error()).
		Block(
			jen.Id("u").Op(",").Id("err").Op(":=").Id("UnmarshalUUID").Call(jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Op("&").Id("u"), jen.Nil()),
		)
}

// URL helpers

func (b *cborHelpersBuilder) buildMarshalURL() jen.Code {
	return jen.Func().Id("MarshalURL").
		Params(jen.Id("u").Qual("net/url", "URL")).
		String().
		Block(
			jen.Return(jen.Id("u").Dot("String").Call()),
		)
}

func (b *cborHelpersBuilder) buildMarshalURLPtr() jen.Code {
	return jen.Func().Id("MarshalURLPtr").
		Params(jen.Id("u").Op("*").Qual("net/url", "URL")).
		Op("*").String().
		Block(
			jen.If(jen.Id("u").Op("==").Nil()).Block(
				jen.Return(jen.Nil()),
			),
			jen.Id("s").Op(":=").Id("u").Dot("String").Call(),
			jen.Return(jen.Op("&").Id("s")),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalURL() jen.Code {
	return jen.Func().Id("UnmarshalURL").
		Params(jen.Id("s").String()).
		Params(jen.Qual("net/url", "URL"), jen.Error()).
		Block(
			jen.Id("u").Op(",").Id("err").Op(":=").Qual("net/url", "Parse").Call(jen.Id("s")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("net/url", "URL").Values(), jen.Err()),
			),
			jen.Return(jen.Op("*").Id("u"), jen.Nil()),
		)
}

func (b *cborHelpersBuilder) buildUnMarshalURLPtr() jen.Code {
	return jen.Func().Id("UnmarshalURLPtr").
		Params(jen.Id("s").Op("*").String()).
		Params(jen.Op("*").Qual("net/url", "URL"), jen.Error()).
		Block(
			jen.If(jen.Id("s").Op("==").Nil()).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			),
			jen.Id("u").Op(",").Id("err").Op(":=").Qual("net/url", "Parse").Call(jen.Op("*").Id("s")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Id("u"), jen.Nil()),
		)
}
