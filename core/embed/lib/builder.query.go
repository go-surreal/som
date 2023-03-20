//go:build embed

package lib

import (
	"strconv"
	"strings"
	"time"
)

type context struct {
	varIndex int
	vars     map[string]any
}

func (c *context) asVar(val any) string {
	index := strconv.Itoa(c.varIndex)
	c.vars[index] = val
	c.varIndex++
	return "$" + index
}

type Query[T any] struct {
	context
	node       string
	fields     string
	groupBy    string
	Where      []Filter[T]
	Sort       []*SortBuilder
	SortRandom bool
	Fetch      []string
	Offset     int
	Limit      int
	Timeout    time.Duration
	Parallel   bool
}

func NewQuery[T any](node string) Query[T] {
	return Query[T]{
		context: context{
			varIndex: 0,
			vars:     map[string]any{},
		},
		node: node,
	}
}

func (q Query[T]) BuildAsAll() *Result {
	q.fields = "*"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsAllIDs() *Result {
	q.fields = "id"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsCount() *Result {
	q.fields = "count()"
	q.groupBy = "id"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) render() string {
	out := "SELECT " + q.fields + " FROM " + q.node + " "

	var t T
	whereStatement := All[T](q.Where).build(&q.context, t)
	if whereStatement != "" {
		out += "WHERE " + whereStatement + " "
	}

	if q.groupBy != "" {
		out += "GROUP BY " + q.groupBy + " "
	}

	if q.SortRandom {
		out += "ORDER BY RAND() "
	} else if len(q.Sort) > 0 {
		var sorts []string
		for _, s := range q.Sort {
			sorts = append(sorts, s.render())
		}
		out += "ORDER BY " + strings.Join(sorts, ", ") + " "
	}

	// LIMIT must come before START.
	if q.Limit > 0 {
		out += "LIMIT " + strconv.Itoa(q.Limit) + " "
	}

	// START must come after LIMIT.
	if q.Offset > 0 {
		out += "START " + strconv.Itoa(q.Offset) + " "
	}

	if len(q.Fetch) > 0 {
		out += "FETCH " + strings.Join(q.Fetch, ", ") + " "
	}

	if q.Timeout > 0 {
		out += "TIMEOUT " + q.Timeout.Round(time.Second).String() + " "
	}

	if q.Parallel {
		out += "PARALLEL"
	}

	return out
}

type Result struct {
	Statement string
	Variables map[string]any
}

type Operator string

const (
	OpEqual            Operator = "="
	OpNotEqual         Operator = "!="
	OpExactlyEqual     Operator = "=="
	OpAnyEqual         Operator = "?="
	OpAllEqual         Operator = "*="
	OpFuzzyMatch       Operator = "~"
	OpFuzzyNotMatch    Operator = "!~"
	OpAnyFuzzyMatch    Operator = "?~"
	OpAllFuzzyMatch    Operator = "*~"
	OpLessThan         Operator = "<"
	OpLessThanEqual    Operator = "<="
	OpGreaterThan      Operator = ">"
	OpGreaterThanEqual Operator = ">="
	OpAdd              Operator = "+"
	OpSub              Operator = "-"
	OpMul              Operator = "*"
	OpDiv              Operator = "/"
	OpAnd              Operator = "AND" // "&&"
	OpOr               Operator = "OR"  // "||"
	OpContains         Operator = "CONTAINS"
	OpContainsNot      Operator = "CONTAINSNOT"
	OpContainsAll      Operator = "CONTAINSALL"
	OpContainsAny      Operator = "CONTAINSANY"
	OpContainsNone     Operator = "CONTAINSNONE"
	OpInside           Operator = "INSIDE"
	OpNotInside        Operator = "NOTINSIDE"
	OpAllInside        Operator = "ALLINSIDE"
	OpAnyInside        Operator = "ANYINSIDE"
	OpNoneInside       Operator = "NONEINSIDE"
	OpOutside          Operator = "OUTSIDE"
	OpIntersect        Operator = "INTERSECTS"
)
