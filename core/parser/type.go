package parser

type Node struct {
	Name           string
	Fields         []Field
	IDType         IDType
	IDEmbed        string
	ComplexID      *FieldComplexID
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
	// Name is an optional index name from `index=<name>` or `unique=<name>`.
	// For regular indexes, this becomes the SurrealDB index name.
	// For unique indexes with a name, fields sharing the same name are
	// grouped into a single composite unique index.
	// If empty, the index name is auto-generated from table and field names.
	Name string

	// Unique indicates this is a unique index.
	Unique bool
}

// SearchInfo holds fulltext search configuration parsed from struct tags.
type SearchInfo struct {
	// ConfigName references a search configuration defined in a //go:build som file.
	ConfigName string
}
