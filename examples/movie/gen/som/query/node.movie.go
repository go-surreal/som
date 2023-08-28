// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package query

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/movie/gen/som/conv"
	lib "github.com/marcbinz/som/examples/movie/gen/som/internal/lib"
	with "github.com/marcbinz/som/examples/movie/gen/som/with"
	model "github.com/marcbinz/som/examples/movie/model"
	marshal "github.com/surrealdb/surrealdb.go/pkg/marshal"
	"strings"
	"time"
)

type Movie struct {
	db    Database
	query lib.Query[model.Movie]
}

func NewMovie(db Database) Movie {
	return Movie{
		db:    db,
		query: lib.NewQuery[model.Movie]("movie"),
	}
}

// Filter adds a where statement to the query to
// select records based on the given conditions.
//
// Use where.All to chain multiple conditions
// together that all need to match.
// Use where.Any to chain multiple conditions
// together where at least one needs to match.
func (q Movie) Filter(filters ...lib.Filter[model.Movie]) Movie {
	q.query.Where = append(q.query.Where, filters...)
	return q
}

// Order sorts the returned records based on the given conditions.
// If multiple conditions are given, they are applied one after the other.
// Note: If OrderRandom is used within the same query,
// it would override the sort conditions.
func (q Movie) Order(by ...*lib.Sort[model.Movie]) Movie {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*lib.SortBuilder)(s))
	}
	return q
}

// OrderRandom sorts the returned records in a random order.
// Note: OrderRandom takes precedence over Order.
func (q Movie) OrderRandom() Movie {
	q.query.SortRandom = true
	return q
}

// Offset skips the first x records for the result set.
func (q Movie) Offset(offset int) Movie {
	q.query.Offset = offset
	return q
}

// Limit restricts the query to return at most x records.
func (q Movie) Limit(limit int) Movie {
	q.query.Limit = limit
	return q
}

// Fetch can be used to return related records.
// This works for both records links and edges.
func (q Movie) Fetch(fetch ...with.Fetch_[model.Movie]) Movie {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
	return q
}

// Timeout adds an execution time limit to the query.
// When exceeded, the query call will return with an error.
func (q Movie) Timeout(timeout time.Duration) Movie {
	q.query.Timeout = timeout
	return q
}

// Parallel tells SurrealDB that individual parts
// of the query can be calculated in parallel.
// This could lead to a faster execution.
func (q Movie) Parallel(parallel bool) Movie {
	q.query.Parallel = parallel
	return q
}

// Count returns the size of the result set, in other words the
// number of records matching the conditions of the query.
func (q Movie) Count(ctx context.Context) (int, error) {
	res := q.query.BuildAsCount()
	raw, err := q.db.Query(ctx, res.Statement, res.Variables)
	if err != nil {
		return 0, err
	}
	var rawCount []marshal.RawQuery[countResult]
	err = marshal.UnmarshalRaw(raw, &rawCount)
	if err != nil {
		return 0, fmt.Errorf("could not count records: %w", err)
	}
	if len(rawCount) < 1 || len(rawCount[0].Result) < 1 {
		return 0, nil
	}
	return rawCount[0].Result[0].Count, nil
}

// CountAsync is the asynchronous version of Count.
func (q Movie) CountAsync(ctx context.Context) *asyncResult[int] {
	return async(ctx, q.Count)
}

// Exists returns whether at least one record for the conditions
// of the query exists or not. In other words it returns whether
// the size of the result set is greater than 0.
func (q Movie) Exists(ctx context.Context) (bool, error) {
	count, err := q.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsAsync is the asynchronous version of Exists.
func (q Movie) ExistsAsync(ctx context.Context) *asyncResult[bool] {
	return async(ctx, q.Exists)
}

// All returns all records matching the conditions of the query.
func (q Movie) All(ctx context.Context) ([]*model.Movie, error) {
	res := q.query.BuildAsAll()
	rawNodes, err := marshal.SmartUnmarshal[conv.Movie](q.db.Query(ctx, res.Statement, res.Variables))
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var nodes []*model.Movie
	for _, rawNode := range rawNodes {
		node := conv.ToMovie(rawNode)
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

// AllAsync is the asynchronous version of All.
func (q Movie) AllAsync(ctx context.Context) *asyncResult[[]*model.Movie] {
	return async(ctx, q.All)
}

// AllIDs returns the IDs of all records matching the conditions of the query.
func (q Movie) AllIDs(ctx context.Context) ([]string, error) {
	res := q.query.BuildAsAllIDs()
	rawNodes, err := marshal.SmartUnmarshal[idNode](q.db.Query(ctx, res.Statement, res.Variables))
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var ids []string
	for _, rawNode := range rawNodes {
		ids = append(ids, rawNode.ID)
	}
	return ids, nil
}

// AllIDsAsync is the asynchronous version of AllIDs.
func (q Movie) AllIDsAsync(ctx context.Context) *asyncResult[[]string] {
	return async(ctx, q.AllIDs)
}

// First returns the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q Movie) First(ctx context.Context) (*model.Movie, error) {
	q.query.Limit = 1
	res, err := q.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("empty result")
	}
	return res[0], nil
}

// FirstAsync is the asynchronous version of First.
func (q Movie) FirstAsync(ctx context.Context) *asyncResult[*model.Movie] {
	return async(ctx, q.First)
}

// FirstID returns the ID of the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q Movie) FirstID(ctx context.Context) (string, error) {
	q.query.Limit = 1
	res, err := q.AllIDs(ctx)
	if err != nil {
		return "", err
	}
	if len(res) < 1 {
		return "", errors.New("empty result")
	}
	return res[0], nil
}

// FirstIDAsync is the asynchronous version of FirstID.
func (q Movie) FirstIDAsync(ctx context.Context) *asyncResult[string] {
	return async(ctx, q.FirstID)
}

// Describe returns a string representation of the query.
// While this might be a valid SurrealDB query, it
// should only be used for debugging purposes.
func (q Movie) Describe() string {
	res := q.query.BuildAsAll()
	return strings.TrimSpace(res.Statement)
}
