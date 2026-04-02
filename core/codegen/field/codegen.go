package field

import (
	"github.com/dave/jennifer/jen"
)

type CodeGenFunc func(ctx Context) jen.Code

func (fn CodeGenFunc) Exec(ctx Context) jen.Code {
	if fn == nil {
		return nil
	}

	return fn(ctx)
}

type CodeGenTuple func(ctx Context) (jen.Code, jen.Code)

func (fn CodeGenTuple) Exec(ctx Context) jen.Code {
	if fn == nil {
		return nil
	}

	a, b := fn(ctx)

	if a == nil || b == nil {
		return nil
	}

	return jen.Add(a).Add(b)
}

type CodeGen struct {
	filterDefine CodeGenFunc
	filterInit   CodeGenTuple
	filterFunc   CodeGenFunc
	// filterExtra generates additional code (types, methods) without disabling filterDefine/filterInit.
	// This is used for wrapper types that need both struct fields AND extra type definitions.
	filterExtra CodeGenFunc

	sortDefine CodeGenFunc
	sortInit   CodeGenFunc
	sortFunc   CodeGenFunc

	fieldDefine CodeGenFunc
	fieldInit   CodeGenFunc
	fieldFunc   CodeGenFunc

	cborMarshal   CodeGenFunc
	cborUnmarshal CodeGenFunc

	// selectDecode generates a custom decode function for SelectField.
	// When set, the generated select method provides a decodeFn that
	// handles types requiring conversion (e.g. time.Time, url.URL).
	selectDecode        CodeGenFunc
	selectDistinctDecode CodeGenFunc
}

func (g *CodeGen) FilterDefine(ctx Context) jen.Code {
	if g.filterFunc.Exec(ctx) != nil {
		return nil
	}

	return g.filterDefine.Exec(ctx)
}

func (g *CodeGen) FilterInit(ctx Context) jen.Code {
	if g.filterFunc.Exec(ctx) != nil {
		return nil
	}

	return g.filterInit.Exec(ctx)
}

func (g *CodeGen) FilterFunc(ctx Context) jen.Code {
	return g.filterFunc.Exec(ctx)
}

func (g *CodeGen) FilterExtra(ctx Context) jen.Code {
	return g.filterExtra.Exec(ctx)
}

func (g *CodeGen) SortDefine(ctx Context) jen.Code {
	return g.sortDefine.Exec(ctx)
}

func (g *CodeGen) SortInit(ctx Context) jen.Code {
	return g.sortInit.Exec(ctx)
}

func (g *CodeGen) SortFunc(ctx Context) jen.Code {
	return g.sortFunc.Exec(ctx)
}

func (g *CodeGen) FieldDefine(ctx Context) jen.Code {
	return g.fieldDefine.Exec(ctx)
}

func (g *CodeGen) FieldInit(ctx Context) jen.Code {
	return g.fieldInit.Exec(ctx)
}

func (g *CodeGen) FieldFunc(ctx Context) jen.Code {
	return g.fieldFunc.Exec(ctx)
}

func (g *CodeGen) CBORMarshal(ctx Context) jen.Code {
	return g.cborMarshal.Exec(ctx)
}

func (g *CodeGen) CBORUnmarshal(ctx Context) jen.Code {
	return g.cborUnmarshal.Exec(ctx)
}

func (g *CodeGen) SelectDecode(ctx Context) jen.Code {
	return g.selectDecode.Exec(ctx)
}

func (g *CodeGen) SelectDistinctDecode(ctx Context) jen.Code {
	return g.selectDistinctDecode.Exec(ctx)
}

// SelectConvertible is implemented by field types whose CBOR representation
// doesn't directly unmarshal into the Go type (e.g. time.Time, time.Duration, url.URL).
// It provides the cbor helper package and function name for element-level conversion,
// which is needed by Slice fields to convert each element.
type SelectConvertible interface {
	SelectCBORHelper(ctx Context) (cborPkg, helperName, ptrHelperName string)
}

// selectDecodeWithHelper generates a decodeFn that unmarshals SELECT VALUE results
// as cbor.RawMessage and converts each element using the given cbor helper function.
// The helper must have signature: func([]byte) (T, error).
// selectDistinctDecodeWithHelper is the same but for distinct queries which return
// a nested array: result: [[raw1, raw2, ...]].
func selectDecodeWithHelper(ctx Context, goType jen.Code, cborPkg, helperName string) jen.Code {
	return jen.Func().Params(
		jen.Id("data").Index().Byte(),
	).Params(
		jen.Index().Add(goType), jen.Error(),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("rawResult").Index().Qual(ctx.pkgInternal(), "QueryResult").Types(
			jen.Qual(cborPkg, "RawMessage"),
		)
		g.If(
			jen.Id("err").Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("rawResult")),
			jen.Id("err").Op("!=").Nil(),
		).Block(jen.Return(jen.Nil(), jen.Id("err")))

		g.If(jen.Len(jen.Id("rawResult")).Op("<").Lit(1).Op("||").Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")).Op("<").Lit(1)).Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)

		g.Id("out").Op(":=").Make(jen.Index().Add(goType), jen.Lit(0), jen.Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")))
		g.For(jen.Id("_").Op(",").Id("raw").Op(":=").Range().Id("rawResult").Index(jen.Lit(0)).Dot("Result")).BlockFunc(func(inner *jen.Group) {
			inner.List(jen.Id("v"), jen.Id("err")).Op(":=").Qual(cborPkg, helperName).Call(jen.Id("raw"))
			inner.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("err")),
			)
			inner.Id("out").Op("=").Append(jen.Id("out"), jen.Id("v"))
		})
		g.Return(jen.Id("out"), jen.Nil())
	})
}

// selectSliceDecodeWithHelper generates a decodeFn for slice fields whose elements
// need conversion. It unmarshals each result row as []RawMessage, then converts
// each element using the given helper.
func selectSliceDecodeWithHelper(ctx Context, sliceType, elemType jen.Code, cborPkg, helperName string) jen.Code {
	return jen.Func().Params(
		jen.Id("data").Index().Byte(),
	).Params(
		jen.Index().Add(sliceType), jen.Error(),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("rawResult").Index().Qual(ctx.pkgInternal(), "QueryResult").Types(
			jen.Index().Qual(cborPkg, "RawMessage"),
		)
		g.If(
			jen.Id("err").Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("rawResult")),
			jen.Id("err").Op("!=").Nil(),
		).Block(jen.Return(jen.Nil(), jen.Id("err")))

		g.If(jen.Len(jen.Id("rawResult")).Op("<").Lit(1).Op("||").Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")).Op("<").Lit(1)).Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)

		g.Id("out").Op(":=").Make(jen.Index().Add(sliceType), jen.Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")))
		g.For(jen.Id("i").Op(",").Id("rawElems").Op(":=").Range().Id("rawResult").Index(jen.Lit(0)).Dot("Result")).BlockFunc(func(inner *jen.Group) {
			inner.Id("converted").Op(":=").Make(jen.Index().Add(elemType), jen.Lit(0), jen.Len(jen.Id("rawElems")))
			inner.For(jen.Id("_").Op(",").Id("raw").Op(":=").Range().Id("rawElems")).BlockFunc(func(inner2 *jen.Group) {
				inner2.List(jen.Id("v"), jen.Id("err")).Op(":=").Qual(cborPkg, helperName).Call(jen.Id("raw"))
				inner2.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Id("err")),
				)
				inner2.Id("converted").Op("=").Append(jen.Id("converted"), jen.Id("v"))
			})
			inner.Id("out").Index(jen.Id("i")).Op("=").Id("converted")
		})
		g.Return(jen.Id("out"), jen.Nil())
	})
}

// selectSliceDistinctDecodeWithHelper is the same but for distinct queries.
func selectSliceDistinctDecodeWithHelper(ctx Context, sliceType, elemType jen.Code, cborPkg, helperName string) jen.Code {
	return jen.Func().Params(
		jen.Id("data").Index().Byte(),
	).Params(
		jen.Index().Add(sliceType), jen.Error(),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("rawResult").Index().Qual(ctx.pkgInternal(), "QueryResult").Types(
			jen.Index().Index().Qual(cborPkg, "RawMessage"),
		)
		g.If(
			jen.Id("err").Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("rawResult")),
			jen.Id("err").Op("!=").Nil(),
		).Block(jen.Return(jen.Nil(), jen.Id("err")))

		g.If(jen.Len(jen.Id("rawResult")).Op("<").Lit(1).Op("||").Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")).Op("<").Lit(1)).Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)

		g.Id("out").Op(":=").Make(jen.Index().Add(sliceType), jen.Lit(0), jen.Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result").Index(jen.Lit(0))))
		g.For(jen.Id("_").Op(",").Id("rawElems").Op(":=").Range().Id("rawResult").Index(jen.Lit(0)).Dot("Result").Index(jen.Lit(0))).BlockFunc(func(inner *jen.Group) {
			inner.Id("converted").Op(":=").Make(jen.Index().Add(elemType), jen.Lit(0), jen.Len(jen.Id("rawElems")))
			inner.For(jen.Id("_").Op(",").Id("raw").Op(":=").Range().Id("rawElems")).BlockFunc(func(inner2 *jen.Group) {
				inner2.List(jen.Id("v"), jen.Id("err")).Op(":=").Qual(cborPkg, helperName).Call(jen.Id("raw"))
				inner2.If(jen.Id("err").Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Id("err")),
				)
				inner2.Id("converted").Op("=").Append(jen.Id("converted"), jen.Id("v"))
			})
			inner.Id("out").Op("=").Append(jen.Id("out"), jen.Id("converted"))
		})
		g.Return(jen.Id("out"), jen.Nil())
	})
}

func selectDistinctDecodeWithHelper(ctx Context, goType jen.Code, cborPkg, helperName string) jen.Code {
	return jen.Func().Params(
		jen.Id("data").Index().Byte(),
	).Params(
		jen.Index().Add(goType), jen.Error(),
	).BlockFunc(func(g *jen.Group) {
		g.Var().Id("rawResult").Index().Qual(ctx.pkgInternal(), "QueryResult").Types(
			jen.Index().Qual(cborPkg, "RawMessage"),
		)
		g.If(
			jen.Id("err").Op(":=").Qual(cborPkg, "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("rawResult")),
			jen.Id("err").Op("!=").Nil(),
		).Block(jen.Return(jen.Nil(), jen.Id("err")))

		g.If(jen.Len(jen.Id("rawResult")).Op("<").Lit(1).Op("||").Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result")).Op("<").Lit(1)).Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)

		g.Id("out").Op(":=").Make(jen.Index().Add(goType), jen.Lit(0), jen.Len(jen.Id("rawResult").Index(jen.Lit(0)).Dot("Result").Index(jen.Lit(0))))
		g.For(jen.Id("_").Op(",").Id("raw").Op(":=").Range().Id("rawResult").Index(jen.Lit(0)).Dot("Result").Index(jen.Lit(0))).BlockFunc(func(inner *jen.Group) {
			inner.List(jen.Id("v"), jen.Id("err")).Op(":=").Qual(cborPkg, helperName).Call(jen.Id("raw"))
			inner.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("err")),
			)
			inner.Id("out").Op("=").Append(jen.Id("out"), jen.Id("v"))
		})
		g.Return(jen.Id("out"), jen.Nil())
	})
}
