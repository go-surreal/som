package som

import (
	"time"
)

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
