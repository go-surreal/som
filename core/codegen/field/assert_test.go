package field

import (
	"testing"

	"github.com/go-surreal/som/core/parser"
	"gotest.tools/v3/assert"
)

// assertSource overrides the constraints the rendering reads, keeping the
// embedded field for everything else. The parser sets them while parsing,
// which is out of reach here.
type assertSource struct {
	parser.Field

	asserts []parser.AssertInfo
}

func (s assertSource) Asserts() []parser.AssertInfo {
	return s.asserts
}

// newAssertField builds a string field carrying the given constraints, with a
// single "phone" config available for reference.
func newAssertField(name string, asserts ...parser.AssertInfo) *baseField {
	source := assertSource{
		Field:   parser.NewFieldString(name),
		asserts: asserts,
	}

	conf := &BuildConfig{
		ToDatabaseName: func(base string) string {
			return base
		},
		AssertConfig: func(configName string) *parser.AssertDef {
			if configName != "phone" {
				return nil
			}
			return &parser.AssertDef{
				Name:    "phone",
				Regex:   `^[0-9+\-]+$`,
				Message: "must be a phone number",
			}
		},
	}

	return &baseField{BuildConfig: conf, source: source, lenFunc: "string::len"}
}

func TestDefineFieldWithoutAsserts(t *testing.T) {
	field := newAssertField("name")

	assert.Equal(t,
		"DEFINE FIELD OVERWRITE name ON TABLE user TYPE string;",
		field.defineField(fieldDef{Name: "name", Table: "user", Type: "string"}),
	)
}

func TestDefineFieldKeepsBuiltinAssert(t *testing.T) {
	field := newAssertField("age")

	assert.Equal(t,
		"DEFINE FIELD OVERWRITE age ON TABLE user TYPE int ASSERT $value >= 0;",
		field.defineField(fieldDef{Name: "age", Table: "user", Type: "int", Assert: "$value >= 0"}),
	)
}

func TestDefineFieldLenRange(t *testing.T) {
	field := newAssertField("name", parser.AssertInfo{Kind: parser.AssertLen, Min: "3", Max: "64"})

	assert.Equal(t,
		`DEFINE FIELD OVERWRITE name ON TABLE user TYPE string ASSERT { `+
			`IF !(string::len($value) >= 3 AND string::len($value) <= 64) `+
			`{ THROW "som_assert:name:length must be between 3 and 64" }; RETURN true; };`,
		field.defineField(fieldDef{Name: "name", Table: "user", Type: "string"}),
	)
}

func TestDefineFieldComposesWithBuiltinAssert(t *testing.T) {
	field := newAssertField("age", parser.AssertInfo{Kind: parser.AssertNum, Min: "0", Max: "130"})

	assert.Equal(t,
		`DEFINE FIELD OVERWRITE age ON TABLE user TYPE int ASSERT { `+
			`IF !($value >= 0 AND $value <= 130) `+
			`{ THROW "som_assert:age:must be between 0 and 130" }; RETURN $value >= -128 AND $value <= 127; };`,
		field.defineField(fieldDef{Name: "age", Table: "user", Type: "int", Assert: "$value >= -128 AND $value <= 127"}),
	)
}

func TestDefineFieldOptionalSkipsAbsentValue(t *testing.T) {
	field := newAssertField("name", parser.AssertInfo{Kind: parser.AssertLen, Min: "3", Max: "3"})

	assert.Equal(t,
		`DEFINE FIELD OVERWRITE name ON TABLE user TYPE option<string> ASSERT { `+
			`IF !($value == NONE OR $value == NULL OR (string::len($value) == 3)) `+
			`{ THROW "som_assert:name:length must be exactly 3" }; RETURN true; };`,
		field.defineField(fieldDef{Name: "name", Table: "user", Type: "option<string>"}),
	)
}

// A slice is stored as an option even when the Go field is not a pointer, so
// its constraint has to tolerate an absent value just the same.
func TestDefineFieldSliceSkipsAbsentValue(t *testing.T) {
	field := newAssertField("tags", parser.AssertInfo{Kind: parser.AssertLen, Max: "5"})
	field.lenFunc = "array::len"

	assert.Equal(t,
		`DEFINE FIELD OVERWRITE tags ON TABLE user TYPE option<array<string>> ASSERT { `+
			`IF !($value == NONE OR $value == NULL OR (array::len($value) <= 5)) `+
			`{ THROW "som_assert:tags:length must be at most 5" }; RETURN true; };`,
		field.defineField(fieldDef{Name: "tags", Table: "user", Type: "option<array<string>>"}),
	)
}

func TestDefineFieldConfigEscapesPattern(t *testing.T) {
	field := newAssertField("phone", parser.AssertInfo{Kind: parser.AssertConfig, ConfigName: "phone"})

	assert.Equal(t,
		`DEFINE FIELD OVERWRITE phone ON TABLE user TYPE string ASSERT { `+
			`IF !(string::matches($value, "^[0-9+\\-]+$")) `+
			`{ THROW "som_assert:phone:must be a phone number" }; RETURN true; };`,
		field.defineField(fieldDef{Name: "phone", Table: "user", Type: "string"}),
	)
}

func TestDefineFieldMultipleConstraints(t *testing.T) {
	field := newAssertField("phone",
		parser.AssertInfo{Kind: parser.AssertLen, Min: "5", Max: ""},
		parser.AssertInfo{Kind: parser.AssertConfig, ConfigName: "phone"},
	)

	assert.Equal(t,
		`DEFINE FIELD OVERWRITE phone ON TABLE user TYPE string ASSERT { `+
			`IF !(string::len($value) >= 5) { THROW "som_assert:phone:length must be at least 5" }; `+
			`IF !(string::matches($value, "^[0-9+\\-]+$")) { THROW "som_assert:phone:must be a phone number" }; `+
			`RETURN true; };`,
		field.defineField(fieldDef{Name: "phone", Table: "user", Type: "string"}),
	)
}
