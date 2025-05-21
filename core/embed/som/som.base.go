//go:build embed

package som

import (
	"github.com/go-surreal/sdbc"
	"time"
)

type ID = sdbc.ID

// type Record = Node // TODO: should we use this to clarify whether a model has edges (node) or not (record)?

type Node struct {
	id *ID
}

func NewNode(id *ID) Node {
	return Node{
		id: id,
	}
}

func (n Node) ID() *ID {
	return n.id
}

// Edge describes an edge between two Node elements.
// It may have its own fields.
type Edge struct {
	id *ID
}

func NewEdge(id *ID) Edge {
	return Edge{
		id: id,
	}
}

func (e Edge) ID() *ID {
	return e.id
}

type Timestamps struct {
	createdAt time.Time
	updatedAt time.Time
}

func NewTimestamps(createdAt *sdbc.DateTime, updatedAt *sdbc.DateTime) Timestamps {
	var ts Timestamps

	if createdAt != nil {
		ts.createdAt = createdAt.Time
	}

	if updatedAt != nil {
		ts.updatedAt = updatedAt.Time
	}

	return ts
}

func (t Timestamps) CreatedAt() time.Time {
	return t.createdAt
}

func (t Timestamps) UpdatedAt() time.Time {
	return t.updatedAt
}

// TODO: implement soft delete feature
// type SoftDelete struct {
// 	deletedAt time.Time
// }
//
// func NewSoftDelete(deletedAt time.Time) SoftDelete {
// 	return SoftDelete{
// 		deletedAt: deletedAt,
// 	}
// }
//
// func (t SoftDelete) DeletedAt() time.Time {
// 	return t.deletedAt
// }

// Enum describes a database type with a fixed set of allowed values.
type Enum string

// Email describes a string field that should contain an email address.
type Email string

// Password describes a string field that should contain a password.
type Password string

// SemVer describes a string field that should contain a semantic version.
type SemVer string

// Password describes a special string field.
// Regarding the generated database query operations, it can only be matched, but never read.
// In a query result, the Password field will always be empty.
// TODO: implement!
// type Password string

// Meta describes a model that is not related to any Node or Edge.
// Instead, it is used to hold metadata that was queried from a Node or Edge.
//
// Applying this struct to a type within your model package will ensure
// that this type is never considered for the generated layer.
// type Meta struct{}

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
// TODO: implement!
// type Info struct {
// 	Time    time.Time
// 	Status  string
// 	Message string
// }

// TODO: below needed?
// type Entity interface {
// 	entity()
// }

// TODO: implement external types
// type External struct {
// 	ID string
// }
