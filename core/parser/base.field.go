package parser

import (
	"fmt"
)

// Field is a single parsed field of a model type.
//
// Every field type is also its own handler: the instance passed to Parse acts
// as the prototype that Match is called on, and Parse returns a fresh instance
// for each matched field. Detection, construction and validation of a field
// therefore all live on one type, in one file.
//
// Field types that som constructs itself (FieldID, FieldComplexID,
// FieldVersion) are never matched against source and implement Field only.
type Field interface {
	fmt.Stringer
	field()

	Name() string
	DBName() string
	Pointer() bool
	Indexes() []IndexInfo
	Search() *SearchInfo
	Asserts() []AssertInfo

	// HasValidate reports whether the type of the field provides a
	// "Validate() error" method that is to be called before a write.
	HasValidate() bool

	// SkipValidate reports whether the field opted out of validation via
	// `som:"novalidate"`, which excludes the field and everything below it.
	SkipValidate() bool

	setName(string)
	setDBName(string)
	setPointer(bool)
	setIndexes([]IndexInfo)
	setSearch(*SearchInfo)
	setAsserts([]AssertInfo)
	setHasValidate(bool)
	setSkipValidate(bool)

	// Validate checks the field on its own, directly after it was parsed.
	Validate() error

	// Verify checks the field against the fully parsed model package, so that
	// it can resolve references to other types. The returned error carries no
	// owning-type prefix, the caller adds that.
	Verify(out *Output) error

	// IsInternal reports whether the field is managed by som itself and may
	// therefore occupy a reserved database name.
	IsInternal() bool
}

// fieldBase carries the state every field shares and provides the default
// behaviour for the parts of Field that most types do not override.
type fieldBase struct {
	name         string
	dbName       string
	pointer      bool
	indexes      []IndexInfo
	search       *SearchInfo
	asserts      []AssertInfo
	hasValidate  bool
	skipValidate bool
}

func newBase(name string) fieldBase {
	return fieldBase{name: name}
}

func (*fieldBase) field() {}

func (f *fieldBase) String() string {
	return f.Name()
}

func (f *fieldBase) Name() string {
	return f.name
}

func (f *fieldBase) setName(name string) {
	f.name = name
}

func (f *fieldBase) DBName() string {
	return f.dbName
}

func (f *fieldBase) setDBName(name string) {
	f.dbName = name
}

func (f *fieldBase) Pointer() bool {
	return f.pointer
}

func (f *fieldBase) setPointer(val bool) {
	f.pointer = val
}

func (f *fieldBase) Indexes() []IndexInfo {
	return f.indexes
}

func (f *fieldBase) setIndexes(indexes []IndexInfo) {
	f.indexes = indexes
}

func (f *fieldBase) Search() *SearchInfo {
	return f.search
}

func (f *fieldBase) setSearch(info *SearchInfo) {
	f.search = info
}

func (f *fieldBase) HasValidate() bool {
	return f.hasValidate
}

func (f *fieldBase) setHasValidate(val bool) {
	f.hasValidate = val
}

func (f *fieldBase) SkipValidate() bool {
	return f.skipValidate
}

func (f *fieldBase) setSkipValidate(val bool) {
	f.skipValidate = val
}

func (f *fieldBase) Asserts() []AssertInfo {
	return f.asserts
}

func (f *fieldBase) setAsserts(asserts []AssertInfo) {
	f.asserts = asserts
}

// assertSupport describes which tag constraints a field type accepts. Named
// configs are rendered as raw SurrealQL and therefore apply to almost any
// field, while len and min/max depend on the database type of the value.
type assertSupport struct {
	length bool
	number bool
	config bool
}

// checkAsserts reports an error for every tag constraint the field type does
// not support, so that a meaningless constraint fails at generation time
// instead of silently never firing.
func (f *fieldBase) checkAsserts(support assertSupport) error {
	for _, assert := range f.asserts {
		switch assert.Kind {
		case AssertLen:
			if !support.length {
				return fmt.Errorf("field %s: len is only supported for string and slice fields", f.name)
			}

		case AssertNum:
			if !support.number {
				return fmt.Errorf("field %s: min and max are only supported for numeric fields", f.name)
			}

		case AssertConfig:
			if !support.config {
				return fmt.Errorf("field %s: assert is not supported for this field type", f.name)
			}
		}
	}

	return nil
}

// rejectAsserts refuses any constraint on a field type that cannot carry one,
// explaining why so the user is not left guessing.
func (f *fieldBase) rejectAsserts(reason string) error {
	if len(f.asserts) > 0 {
		return fmt.Errorf("field %s: %s", f.name, reason)
	}

	return nil
}

// validate runs the checks every field type shares, parameterised with the tag
// constraints the concrete type accepts.
func (f *fieldBase) validate(support assertSupport) error {
	if f.search != nil {
		return fmt.Errorf("field %s: fulltext index only supports string types (string, *string, []string, []*string, *[]string, *[]*string)", f.name)
	}

	return f.checkAsserts(support)
}

func (f *fieldBase) Validate() error {
	return f.validate(assertSupport{config: true})
}

func (f *fieldBase) Verify(out *Output) error {
	if f.search != nil && !searchExists(f.search.ConfigName, out) {
		return fmt.Errorf("field %q references unknown search config %q", f.name, f.search.ConfigName)
	}

	for _, assert := range f.asserts {
		if assert.Kind == AssertConfig && !assertExists(assert.ConfigName, out) {
			return fmt.Errorf("field %q references unknown assert config %q", f.name, assert.ConfigName)
		}
	}

	return nil
}

func (f *fieldBase) IsInternal() bool {
	return false
}
