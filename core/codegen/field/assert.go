package field

import (
	"fmt"
	"strings"

	"github.com/go-surreal/som/core/parser"
)

// assertMarker prefixes the message thrown by a violated constraint, so that
// the generated repository can tell a constraint violation apart from any
// other database error and report it as a som.AssertError.
const assertMarker = "som_assert"

// assertClause is a single constraint, ready to be rendered into the schema.
type assertClause struct {
	// Expr is the SurrealQL expression that must hold for the value.
	Expr string
	// Message explains the violation to the caller.
	Message string
}

// assertClauses renders the constraints declared on the field via struct tags.
// Constraint support is checked while parsing, so an unsupported combination
// reaching this point is a bug in the parser rather than user error.
//
// An optional field is one the database may hold no value for, which is not
// the same as the Go field being a pointer: a slice is stored as an option
// too, so that a nil slice stays distinguishable from an empty one.
func (f *baseField) assertClauses(optional bool) []assertClause {
	var clauses []assertClause

	for _, assert := range f.source.Asserts() {
		var clause assertClause

		switch assert.Kind {
		case parser.AssertLen:
			clause = f.lenClause(assert)

		case parser.AssertNum:
			clause = numClause(assert)

		case parser.AssertConfig:
			clause = f.configClause(assert)

		default:
			panic(fmt.Sprintf("unmapped assert kind: %d", assert.Kind))
		}

		// An absent value cannot violate a constraint, otherwise every
		// optional field would become mandatory.
		if optional {
			clause.Expr = fmt.Sprintf("$value == NONE OR $value == NULL OR (%s)", clause.Expr)
		}

		clauses = append(clauses, clause)
	}

	return clauses
}

// lenClause constrains the length of a string or the item count of a slice.
func (f *baseField) lenClause(assert parser.AssertInfo) assertClause {
	length := f.lenFunc + "($value)"

	switch {
	case assert.Min == assert.Max:
		return assertClause{
			Expr:    fmt.Sprintf("%s == %s", length, assert.Min),
			Message: fmt.Sprintf("length must be exactly %s", assert.Min),
		}

	case assert.Max == "":
		return assertClause{
			Expr:    fmt.Sprintf("%s >= %s", length, assert.Min),
			Message: fmt.Sprintf("length must be at least %s", assert.Min),
		}

	case assert.Min == "":
		return assertClause{
			Expr:    fmt.Sprintf("%s <= %s", length, assert.Max),
			Message: fmt.Sprintf("length must be at most %s", assert.Max),
		}

	default:
		return assertClause{
			Expr:    fmt.Sprintf("%s >= %s AND %s <= %s", length, assert.Min, length, assert.Max),
			Message: fmt.Sprintf("length must be between %s and %s", assert.Min, assert.Max),
		}
	}
}

// numClause constrains the range of a numeric value.
func numClause(assert parser.AssertInfo) assertClause {
	switch {
	case assert.Max == "":
		return assertClause{
			Expr:    fmt.Sprintf("$value >= %s", assert.Min),
			Message: fmt.Sprintf("must be at least %s", assert.Min),
		}

	case assert.Min == "":
		return assertClause{
			Expr:    fmt.Sprintf("$value <= %s", assert.Max),
			Message: fmt.Sprintf("must be at most %s", assert.Max),
		}

	default:
		return assertClause{
			Expr:    fmt.Sprintf("$value >= %s AND $value <= %s", assert.Min, assert.Max),
			Message: fmt.Sprintf("must be between %s and %s", assert.Min, assert.Max),
		}
	}
}

// configClause renders a constraint declared via define.Assert.
func (f *baseField) configClause(assert parser.AssertInfo) assertClause {
	config := f.AssertConfig(assert.ConfigName)
	if config == nil {
		// The parser rejects unknown config names before codegen runs.
		panic(fmt.Sprintf("unknown assert config: %s", assert.ConfigName))
	}

	var exprs []string
	if config.Regex != "" {
		exprs = append(exprs, fmt.Sprintf("string::matches($value, %s)", quoteSurreal(config.Regex)))
	}
	if config.Raw != "" {
		exprs = append(exprs, "("+config.Raw+")")
	}

	message := config.Message
	if message == "" {
		message = fmt.Sprintf("must conform to %s", config.Name)
	}

	return assertClause{
		Expr:    strings.Join(exprs, " AND "),
		Message: message,
	}
}

// assertExpression combines the type's built-in constraint with the ones
// declared via struct tags.
//
// Without tag constraints the built-in assert is emitted as is, keeping the
// schema of an unconstrained field unchanged. Otherwise the constraints become
// a block that throws a marked message per violation, which is what turns a
// generic database error into a message naming the field and the rule it broke.
func (f *baseField) assertExpression(fieldPath, dbType, builtin string) string {
	clauses := f.assertClauses(strings.HasPrefix(dbType, "option<"))
	if len(clauses) == 0 {
		return builtin
	}

	if builtin == "" {
		builtin = "true"
	}

	var block strings.Builder
	block.WriteString("{ ")

	for _, clause := range clauses {
		marker := fmt.Sprintf("%s:%s:%s", assertMarker, fieldPath, clause.Message)
		fmt.Fprintf(&block, "IF !(%s) { THROW %s }; ", clause.Expr, quoteSurreal(marker))
	}

	fmt.Fprintf(&block, "RETURN %s; }", builtin)

	return block.String()
}

// quoteSurreal renders a Go string as a SurrealQL string literal.
func quoteSurreal(value string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value)

	return `"` + escaped + `"`
}
