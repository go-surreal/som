package parser

type Node struct {
	Name           string
	Fields         []Field
	IDType         IDType
	IDEmbed        string
	ComplexID      *FieldComplexID
	Unions         []Membership
	Timestamps     bool
	OptimisticLock bool
	Changefeed     string
	SoftDelete     bool
	Expiry         bool
	ExpiryDuration string
}

// Union is a set of node models sharing a common interface in the model
// package. It has no table of its own: a field typed as a union maps to a
// multi-table record link (record<a|b>), and an edge endpoint typed as a
// union maps to TYPE RELATION IN a|b.
type Union struct {
	Name    string
	Members []string

	// Open marks a union declared with som.OpenUnion instead of som.Union. It
	// requires no marker method, so a node can be a member of any number of
	// open unions, at the cost of the compiler no longer checking membership.
	Open bool
}

// Membership is a node's declaration that it belongs to a union. A sealed union
// is joined by embedding som.Part, which lets the compiler check the
// membership. An open one is joined by a blank som.Part field, which declares
// it without adding a method.
type Membership struct {
	Union  string
	Sealed bool
}

type Edge struct {
	Name           string
	In             Field
	Out            Field
	Fields         []Field
	Timestamps     bool
	OptimisticLock bool
	Changefeed     string
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

// Fragment is a projection of a Node: a struct holding a subset of the node's
// fields. It has no table of its own — fragments narrow the select list of a
// query on the parent node's table. Its fields are taken over from the parent
// node, so that the database name and type of a field cannot diverge.
type Fragment struct {
	Name   string
	Parent string
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
