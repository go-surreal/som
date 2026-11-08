package field

import (
	"testing"

	"github.com/go-surreal/som/core/parser"
	"gotest.tools/v3/assert"
)

func newPtrStructSliceField(name string, asserts ...parser.AssertInfo) *Slice {
	field := newAssertField(name, asserts...)
	field.lenFunc = "array::len"

	return &Slice{baseField: field}
}

func TestDefinePtrStructSliceWithoutAsserts(t *testing.T) {
	field := newPtrStructSliceField("items")

	assert.Equal(t,
		`IF (INFO FOR TABLE user STRUCTURE).fields[WHERE name = "meta.items"][0].kind = "none | array<none | object>" `+
			`{ ALTER FIELD meta.items ON TABLE user DROP ASSERT; } `+
			`ELSE { DEFINE FIELD OVERWRITE meta.items ON TABLE user TYPE option<array<option<object>>>; };`,
		field.definePtrStructSlice("user", "meta."),
	)
}

func TestDefinePtrStructSliceAltersAssert(t *testing.T) {
	field := newPtrStructSliceField("items", parser.AssertInfo{Kind: parser.AssertLen, Max: "5"})

	assertBlock := `{ IF !($value == NONE OR $value == NULL OR (array::len($value) <= 5)) ` +
		`{ THROW "som_assert:items:length must be at most 5" }; RETURN true; }`

	assert.Equal(t,
		`IF (INFO FOR TABLE user STRUCTURE).fields[WHERE name = "items"][0].kind = "none | array<none | object>" `+
			`{ ALTER FIELD items ON TABLE user ASSERT `+assertBlock+`; } `+
			`ELSE { DEFINE FIELD OVERWRITE items ON TABLE user TYPE option<array<option<object>>> ASSERT `+assertBlock+`; };`,
		field.definePtrStructSlice("user", ""),
	)
}
