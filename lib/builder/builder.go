package builder

import (
	"strconv"
	"strings"
	"time"
)

type context struct {
	varIndex rune
	vars     map[string]any
}

func (c *context) asVar(val any) string {
	index := string(c.varIndex)
	c.vars[index] = val
	c.varIndex++
	return "$" + index
}

type Query struct {
	*context
	node       string
	fields     string
	groupBy    string
	Where      []Where
	Sort       []*Sort
	SortRandom bool
	Fetch      []string
	Offset     int
	Limit      int
	Timeout    time.Duration
	Parallel   bool
}

func NewQuery(node string) *Query {
	return &Query{
		context: &context{
			varIndex: 'A',
			vars:     map[string]any{},
		},
		node: node,
	}
}

func (q Query) BuildAsAll() *Result {
	q.fields = "*"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query) BuildAsAllIDs() *Result {
	q.fields = "id"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query) BuildAsCount() *Result {
	q.fields = "count()"
	q.groupBy = "id"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query) render() string {
	out := "SELECT " + q.fields + " FROM " + q.node + " "

	whereStatement := WhereAll{Where: q.Where}.render(q.context)
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
