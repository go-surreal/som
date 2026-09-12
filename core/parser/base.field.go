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

	setName(string)
	setDBName(string)
	setPointer(bool)
	setIndexes([]IndexInfo)
	setSearch(*SearchInfo)

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
	name    string
	dbName  string
	pointer bool
	indexes []IndexInfo
	search  *SearchInfo
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

func (f *fieldBase) Validate() error {
	if f.search != nil {
		return fmt.Errorf("field %s: fulltext index only supports string types (string, *string, []string, []*string, *[]string, *[]*string)", f.name)
	}
	return nil
}

func (f *fieldBase) Verify(out *Output) error {
	if f.search != nil && !searchExists(f.search.ConfigName, out) {
		return fmt.Errorf("field %q references unknown search config %q", f.name, f.search.ConfigName)
	}
	return nil
}

func (f *fieldBase) IsInternal() bool {
	return false
}
