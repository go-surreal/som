package parser

type Node struct {
	Name           string
	Fields         []Field
	IDType         IDType
	IDEmbed        string
	Timestamps     bool
	OptimisticLock bool
	SoftDelete     bool
}

type Edge struct {
	Name           string
	In             Field
	Out            Field
	Fields         []Field
	Timestamps     bool
	OptimisticLock bool
	SoftDelete     bool
}

type Struct struct {
	Name   string
	Fields []Field
}

type Enum struct {
	Name string
}

type EnumValue struct {
	Enum     string
	Variable string
	Value    string
}

// IndexInfo holds index configuration parsed from struct tags.
type IndexInfo struct {
	// Name is the index name. If empty, auto-generated from table and field names.
	Name string

	// Unique indicates this is a unique index.
	Unique bool

	// UniqueName is the composite unique index identifier from `unique(name)`.
	// Fields with the same UniqueName are grouped into a single composite unique index.
	UniqueName string
}

// SearchInfo holds fulltext search configuration parsed from struct tags.
type SearchInfo struct {
	// ConfigName references a search configuration defined in a //go:build som file.
	ConfigName string
}
