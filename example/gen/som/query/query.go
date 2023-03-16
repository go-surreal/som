// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package query

import (
	"context"
	"time"
	
	"github.com/marcbinz/som"
	"github.com/marcbinz/som/example/gen/som/with"
	"github.com/marcbinz/som/lib"
)

type Database interface {
	Query(statement string, vars any) (any, error)
}

type idNode struct {
	ID string
}

type countResult struct {
	Count int
}

type Query[M som.Node] interface{

	// Filter adds a where statement to the query to
	// select records based on the given conditions.
	//
	// Use where.All to chain multiple conditions
	// together that all need to match.
	// Use where.Any to chain multiple conditions
	// together where at least one needs to match.
	Filter(filters ...lib.Filter[M]) Query[M]

	// Order sorts the returned records based on the given conditions.
	// If multiple conditions are given, they are applied one after the other.
	// Note: If OrderRandom is used within the same query,
	// it would override the sort conditions.
	Order(by ...*lib.Sort[M]) Query[M]

	// OrderRandom sorts the returned records in a random order.
	// Note: OrderRandom takes precedence over Order.
	OrderRandom() Query[M]

	// Offset skips the first x records for the result set.
	Offset(offset int) Query[M]

	// Limit restricts the query to return at most x records.
	Limit(limit int) Query[M]

	// Fetch can be used to return related records.
	// This works for both records links and edges.
	Fetch(fetch ...with.Fetch_[M]) Query[M]

	// Timeout adds an execution time limit to the query.
	// When exceeded, the query call will return with an error.
	Timeout(timeout time.Duration) Query[M]

	// Parallel tells SurrealDB that individual parts
	// of the query can be calculated in parallel.
	// This could lead to a faster execution.
	Parallel(parallel bool) Query[M]

	// Count returns the size of the result set, in other words the
	// number of records matching the conditions of the query.
	Count(ctx context.Context) (int, error)

	// Exists returns whether at least one record for the conditons
	// of the query exists or not. In other words it returns whether
	// the size of the result set is greater than 0.
	Exists(ctx context.Context) (bool, error)

	// All returns all records matching the conditions of the query.
	All(ctx context.Context) ([]*M, error)

	// AllIDs returns the IDs of all records matching the conditions of the query.
	AllIDs(ctx context.Context) ([]string, error)

	// First returns the first record matching the conditions of the query.
	// This comes in handy when using a filter for a field with unique values or when
	// sorting the result set in a specific order where only the first result is relevant.
	First(ctx context.Context) (*M, error)

	// FirstID returns the ID of the first record matching the conditions of the query.
	// This comes in handy when using a filter for a field with unique values or when
	// sorting the result set in a specific order where only the first result is relevant.
	FirstID(ctx context.Context) (string, error)

	// Describe returns a string representation of the query.
	// While this might be a valid SurrealDB query, it
	// should only be used for debugging purposes.
	Describe() string
}