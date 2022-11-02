package builder

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
	node  string
	Where []Where
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

func (q Query) render() string {
	whereStatement := WhereAll{Where: q.Where}.render(q.context)

	where := ""
	if whereStatement != "" {
		where = " WHERE " + whereStatement
	}

	return "SELECT * FROM " + q.node + where
}

func Build(q *Query) *Result {
	return &Result{
		Statement: q.render(),
		Variables: q.context.vars,
	}
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
