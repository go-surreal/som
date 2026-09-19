//go:build embed

package define

// AssertBuilder builds a named field constraint that model fields reference
// via the `som:"assert=<name>"` struct tag.
//
// Constraints that fit into a struct tag without quoting or escaping (len, min,
// max) are declared on the field itself. Everything else is declared here, so
// that patterns and expressions stay real Go string literals.
type AssertBuilder struct {
	name    string
	regex   string
	raw     string
	message string
}

// Assert creates a new named field constraint.
func Assert(name string) *AssertBuilder {
	return &AssertBuilder{name: name}
}

// Regex requires the field value to match the given pattern.
//
// The pattern is evaluated by SurrealDB, not by Go's regexp package. Both are
// RE2-based, but only the subset both engines agree on is portable.
func (b *AssertBuilder) Regex(pattern string) *AssertBuilder {
	b.regex = pattern
	return b
}

// Raw requires the given SurrealQL expression to hold for the field value.
// The value under validation is available as $value.
//
// This is the escape hatch for constraints the builder does not model. The
// expression is embedded into the schema verbatim and is not validated by som.
func (b *AssertBuilder) Raw(expr string) *AssertBuilder {
	b.raw = expr
	return b
}

// Message sets the error message reported when the constraint is violated.
// If unset, a message is derived from the constraint itself.
func (b *AssertBuilder) Message(msg string) *AssertBuilder {
	b.message = msg
	return b
}

// assertJSON is the JSON representation of a named field constraint.
type assertJSON struct {
	Name    string `json:"name"`
	Regex   string `json:"regex,omitempty"`
	Raw     string `json:"raw,omitempty"`
	Message string `json:"message,omitempty"`
}

// toJSON converts the assert builder to its JSON representation.
func (b *AssertBuilder) toJSON() assertJSON {
	return assertJSON{
		Name:    b.name,
		Regex:   b.regex,
		Raw:     b.raw,
		Message: b.message,
	}
}
