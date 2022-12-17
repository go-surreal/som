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

// Info holds information about a single database operation.
// It is used as target to hold said information when building
// an operation using the WithInfo() method. The generated
// builder for each model provides this capability.
//
// Example:
// Take a model named "Model" for which the som code is generated.
// Accessing the database operation happens as usual for example via
// client.Model().Create() or client.Model().Query(). Extracting the
// Info out of those operations is as simple as:
//
// var info *som.Info
// client.Model().WithInfo(info).Create()
// client.Model().WithInfo(info).Query()
//
// Please note: When using the same base for multiple operations, the Info
// struct will only ever hold the information of the last operation.
type Info struct {
	Time    time.Time
	Status  string
	Message string
}

type Entity interface {
	entity()
}
