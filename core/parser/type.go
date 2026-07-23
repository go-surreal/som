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
	TTL            bool
	TTLExpiry      string
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

// View is a read-only, pre-computed table view. It has a struct shape
// (its projected columns) like a Node, but no ID type, features or
// write operations. The SELECT that populates it is supplied separately
// via a //go:build som definition and linked back by database name.
type View struct {
	Name   string
	Fields []Field
}

// Sink is a write-only ingestion table backed by a DEFINE TABLE ... DROP
// statement. It has a struct shape (its columns) like a Node but no ID,
// no features and only create operations; rows are discarded immediately
// after write, so they cannot be read back or linked to.
type Sink struct {
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
