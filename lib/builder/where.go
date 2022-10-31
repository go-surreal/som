package builder

import (
	"strings"
)

type Where interface {
	Block
	where()
}

type Predicate struct {
	Op      Operator
	Key     string
	Val     any
	IsCount bool
}

func (p Predicate) render(ctx *context) string {
	res := p.Key + " " + string(p.Op) + " " + ctx.asVar(p.Val)
	if p.IsCount {
		res = "COUNT " + res
	}
	return res
}

func (p Predicate) where() {}

type WhereAll struct {
	Where []Where
}

func (p WhereAll) render(ctx *context) string {
	var parts []string

	for _, where := range p.Where {
		parts = append(parts, where.render(ctx))
	}

	return strings.Join(parts, " "+string(OpAnd)+" ")
}

func (p WhereAll) where() {}

type WhereAny struct {
	Where []Where
}

func (p WhereAny) render(ctx *context) string {
	var parts []string

	for _, where := range p.Where {
		parts = append(parts, where.render(ctx))
	}

	return strings.Join(parts, " "+string(OpOr)+" ")
}

func (p WhereAny) where() {}
