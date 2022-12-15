package som

import (
	"time"
)

type Record = Node // TODO: should we use this?

type Node struct {
	// include query info into each node resulting from a query?:
	// status string
	// time   string
	// extract via som.Info(someNode) -> som.Info ?
}

// Edge describes an edge between two Node elements.
// It may have its own fields.
type Edge struct{}

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Enum describes a database type with a fixed set of allowed values.
type Enum string

// Password describes a special string field.
// Regarding the generated database query operations, it can only be matched, but never read.
// In a query result, the Password field will always be empty.
type Password string

// Meta describes a model that is not related to any Node or Edge.
// Instead, it is used to hold metadata that was queried from a Node or Edge.
//
// Applying this struct to a type within your model package will ensure
// that this type is never considered for the generated layer.
type Meta struct{}

// External describes a database entity that only holds a reference
// to an external object. An example for an external object would be
// some data that is fetched from a remote API.
//
// The fields of an External model will not be stored in the database.
// Furthermore, the models can only be read and assigned as links or edges
// to other nodes, but not created or updates. This should happen directly
// via the remote API if needed. Primary use case for the External type are
// referencing purposes, hence the read-only approach.
type External struct{}
