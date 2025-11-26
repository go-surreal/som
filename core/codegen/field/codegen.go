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

	sortDefine CodeGenFunc
	sortInit   CodeGenFunc
	sortFunc   CodeGenFunc

	cborMarshal   CodeGenFunc
	cborUnmarshal CodeGenFunc
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

func (g *CodeGen) SortDefine(ctx Context) jen.Code {
	return g.sortDefine.Exec(ctx)
}

func (g *CodeGen) SortInit(ctx Context) jen.Code {
	return g.sortInit.Exec(ctx)
}

func (g *CodeGen) SortFunc(ctx Context) jen.Code {
	return g.sortFunc.Exec(ctx)
}

func (g *CodeGen) CBORMarshal(ctx Context) jen.Code {
	return g.cborMarshal.Exec(ctx)
}

func (g *CodeGen) CBORUnmarshal(ctx Context) jen.Code {
	return g.cborUnmarshal.Exec(ctx)
}
