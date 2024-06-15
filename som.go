package som

import (
	"time"
)

// type Record = Node // TODO: should we use this to clarify whether a model has edges (node) or not (record)?

// type Record[T any] struct {
// 	id T
// }
//
// func NewRecord[T any](id T) Record[T] {
// 	return Record[T]{
// 		id: id,
// 	}
// }
//
// func (r Record[T]) ID() T {
// 	return r.id
// }
//
// type TimeSeries Record[TimeSeriesID]
//
// type TimeSeriesID struct {
// 	timestamp time.Time
// }

type Node struct {
	id string
}

func NewNode(id string) Node {
	return Node{
		id: id,
	}
}

func (n Node) ID() string {
	return n.id
}

// Edge describes an edge between two Node elements.
// It may have its own fields.
type Edge[I, O any] struct {
	id string
}

func NewEdge[I, O any](id string) Edge[I, O] {
	return Edge[I, O]{
		id: id,
	}
}

func (e Edge[I, O]) ID() string {
	return e.id
}

type Timestamps struct {
	createdAt time.Time
	updatedAt time.Time
}

func NewTimestamps(createdAt *time.Time, updatedAt *time.Time) Timestamps {
	var ts Timestamps

	if createdAt != nil {
		ts.createdAt = *createdAt
	}

	if updatedAt != nil {
		ts.updatedAt = *updatedAt
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
