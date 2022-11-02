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
	if p.IsCount {
		return "count(" + p.Key + ") " + string(p.Op) + " " + ctx.asVar(p.Val)
	}
	return p.Key + " " + string(p.Op) + " " + ctx.asVar(p.Val)
}

func (p Predicate) where() {}

type WhereAll struct {
	Where []Where
}

func (p WhereAll) render(ctx *context) string {
	if len(p.Where) < 1 {
		return ""
	}

	var parts []string
	for _, where := range p.Where {
		if part := where.render(ctx); part != "" {
			parts = append(parts, where.render(ctx))
		}
	}

	if len(parts) < 1 {
		return ""
	}

	return "(" + strings.Join(parts, " "+string(OpAnd)+" ") + ")"
}

func (p WhereAll) where() {}

type WhereAny struct {
	Where []Where
}

func (p WhereAny) render(ctx *context) string {
	if len(p.Where) < 1 {
		return ""
	}

	var parts []string
	for _, where := range p.Where {
		if part := where.render(ctx); part != "" {
			parts = append(parts, where.render(ctx))
		}
	}

	if len(parts) < 1 {
		return ""
	}

	return "(" + strings.Join(parts, " "+string(OpOr)+" ") + ")"
}

func (p WhereAny) where() {}
