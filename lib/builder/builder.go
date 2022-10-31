package builder

import (
	"fmt"
)

type context struct {
	varIndex rune
	vars     map[rune]any
}

func (c *context) asVar(val any) string {
	index := c.varIndex
	c.vars[index] = val
	c.varIndex++
	return "$" + string(index)
}

type Query struct {
	*context
	Where []Block
}

func NewQuery() *Query {
	return &Query{
		context: &context{
			varIndex: 'A',
			vars:     map[rune]any{},
		},
	}
}

func (q Query) render() string {
	return "SELECT * FROM <table>"
}

func Build(q *Query) {
	for _, where := range q.Where {
		fmt.Println(where.render(q.context))
	}

	fmt.Println("\nVariables:")

	for key, value := range q.context.vars {
		fmt.Println(string(key), ":", value)
	}
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
