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

func (c *context) Vars() map[string]any {
	return c.vars
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
	live       bool
	fields     string
	groupBy    string
	groupAll   bool
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
	q.groupAll = true

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsLive() *Result {
	q.live = true
	q.fields = "*"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) BuildAsLiveDiff() *Result {
	q.live = true
	q.fields = "DIFF"

	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
}

func (q Query[T]) render() string {
	var out strings.Builder

	out.WriteString(strings.Join([]string{"SELECT", q.fields, "FROM", q.node}, " "))

	var t T
	whereStatement := All[T](q.Where).build(&q.context, t)
	if whereStatement != "" {
		out.WriteString(" WHERE ")
		out.WriteString(whereStatement)
	}

	if !q.live && q.groupBy != "" {
		out.WriteString(" GROUP BY ")
		out.WriteString(q.groupBy)
	}

	if !q.live && q.groupAll {
		out.WriteString(" GROUP ALL")
	}

	if !q.live && q.SortRandom {
		out.WriteString(" ORDER BY RAND()")
	} else if !q.live && len(q.Sort) > 0 {
		var sorts []string
		for _, s := range q.Sort {
			sorts = append(sorts, s.render())
		}

		out.WriteString(" ORDER BY ")
		out.WriteString(strings.Join(sorts, ", "))
	}

	// LIMIT must come before START.
	if !q.live && q.Limit > 0 {
		out.WriteString(" LIMIT ")
		out.WriteString(strconv.Itoa(q.Limit))
	}

	// START must come after LIMIT.
	if !q.live && q.Offset > 0 {
		out.WriteString(" START ")
		out.WriteString(strconv.Itoa(q.Offset))
	}

	if len(q.Fetch) > 0 {
		out.WriteString(" FETCH ")
		out.WriteString(strings.Join(q.Fetch, ", "))
	}

	if !q.live && q.Timeout > 0 {
		out.WriteString(" TIMEOUT ")
		out.WriteString(q.Timeout.Round(time.Second).String())
	}

	if !q.live && q.Parallel {
		out.WriteString(" PARALLEL")
	}

	return out.String()
}

type Result struct {
	Statement string
	Variables map[string]any
}

type Operator string

const (
	OpAnd Operator = "AND" // ("&&") Checks whether both of two values are truthy.
	OpOr  Operator = "OR"  // ("||") Checks whether either of two values is truthy.

	OpEqual         Operator = "="  // ("IS") Check whether two values are equal.
	OpNotEqual      Operator = "!=" // ("IS NOT") Check whether two values are not equal.
	OpExactlyEqual  Operator = "==" // Check whether two values are exactly equal.
	OpFuzzyMatch    Operator = "~"  // Compare two values for equality using fuzzy matching.
	OpFuzzyNotMatch Operator = "!~" // Compare two values for inequality using fuzzy matching.

	OpAnyEqual      Operator = "?=" // Check whether any value in a set is equal to a value.
	OpAllEqual      Operator = "*=" // Check whether all values in a set are equal to a value.
	OpAnyFuzzyMatch Operator = "?~" // Check whether any value in a set is equal to a value using fuzzy matching.
	OpAllFuzzyMatch Operator = "*~" // Check whether all values in a set are equal to a value using fuzzy matching.

	OpLessThan         Operator = "<"  // Check whether a value is less than another value.
	OpLessThanEqual    Operator = "<=" // Check whether a value is less than or equal to another value.
	OpGreaterThan      Operator = ">"  // Check whether a value is greater than another value.
	OpGreaterThanEqual Operator = ">=" // Check whether a value is greater than or equal to another value.

	OpAdd   Operator = "+"  // 	Add two values together.
	OpSub   Operator = "-"  // Subtract a value from another value.
	OpMul   Operator = "×"  // ("*") Multiply two values together.
	OpDiv   Operator = "÷"  // ("/") Divide a value by another value.
	OpRaise Operator = "**" // Raises a base value by another value.

	OpNot                  Operator = "!"  // Reverses the truthiness of a value.
	OpTruth                Operator = "!!" // Determines the truthiness of a value.
	OpEitherTrueAndNotNull Operator = "??" // Check whether either of two values are truthy and not NULL.
	OpEitherTrue           Operator = "?:" // Check whether either of two values are truthy.

	OpContains     Operator = "∋" // ("CONTAINS") Checks whether a value contains another value.
	OpContainsNot  Operator = "∌" // ("CONTAINSNOT") Checks whether a value does not contain another value.
	OpContainsAll  Operator = "⊇" // ("CONTAINSALL") Checks whether a value contains all other values.
	OpContainsAny  Operator = "⊃" // ("CONTAINSANY") Checks whether a value contains any other value.
	OpContainsNone Operator = "⊅" // ("CONTAINSNONE") Checks whether a value contains none of the following values.

	OpIn     Operator = "∈" // ("INSIDE") Checks whether a value is contained within another value. - TODO: geo!
	OpNotIn  Operator = "∉" // ("NOTINSIDE" | "NOT IN") Checks whether a value is not contained within another value. - TODO: geo!
	OpAllIn  Operator = "⊆" // ("ALLINSIDE") Checks whether all values are contained within other values.
	OpAnyIn  Operator = "⊂" // ("ANYINSIDE") Checks whether any value is contained within other values.
	OpNoneIn Operator = "⊄" // ("NONEINSIDE") Checks whether no value is contained within other values.

	OpGeoOutside    Operator = "OUTSIDE"    // Checks whether a geometry type is outside another geometry type.
	OpGeoIntersects Operator = "INTERSECTS" // Checks whether a geometry type intersects another geometry type.

	OpSearch Operator = "@@" // ("@[ref]@") Checks whether the terms are found in a full-text indexed field. - TODO!

	OpX = "<|4|> or <|3,HAMMING|>" // KNN - TODO!

	CastInt Operator = "<int>"
)
