package predicate

type Predicate interface{}

type User struct{}

const (
	opEqual            = "="
	opNotEqual         = "!="
	opExactlyEqual     = "=="
	opAnyEqual         = "?="
	opAllEqual         = "*="
	opFuzzyMatch       = "~"
	opFuzzyNotMatch    = "!~"
	opAnyFuzzyMatch    = "?~"
	opAllFuzzyMatch    = "*~"
	opLessThan         = "<"
	opLessThanEqual    = "<="
	opGreaterThan      = ">"
	opGreaterThanEqual = ">="
	opAdd              = "+"
	opSub              = "-"
	opMul              = "*"
	opDiv              = "/"
	opAnd              = "AND" // "&&"
	opOr               = "OR"  // "||"
	// CONTAINS, CONTAINSNOT, CONTAINSALL, CONTAINSANY,CONTAINSNONE
	// INSIDE,NOTINSIDE,ALLINSIDE,ANYINSIDE,NONEINSIDE
	// OUTSIDE,INTERSECTS
)
