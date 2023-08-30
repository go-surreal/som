// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package query

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/testing/gen/som/conv"
	lib "github.com/marcbinz/som/examples/testing/gen/som/internal/lib"
	with "github.com/marcbinz/som/examples/testing/gen/som/with"
	model "github.com/marcbinz/som/examples/testing/model"
	"strings"
	"time"
)

type FieldsLikeDBResponse struct {
	db        Database
	query     lib.Query[model.FieldsLikeDBResponse]
	unmarshal func(buf []byte, val any) error
}

func NewFieldsLikeDBResponse(db Database, unmarshal func(buf []byte, val any) error) FieldsLikeDBResponse {
	return FieldsLikeDBResponse{
		db:        db,
		query:     lib.NewQuery[model.FieldsLikeDBResponse]("fields_like_db_response"),
		unmarshal: unmarshal,
	}
}

// Filter adds a where statement to the query to
// select records based on the given conditions.
//
// Use where.All to chain multiple conditions
// together that all need to match.
// Use where.Any to chain multiple conditions
// together where at least one needs to match.
func (q FieldsLikeDBResponse) Filter(filters ...lib.Filter[model.FieldsLikeDBResponse]) FieldsLikeDBResponse {
	q.query.Where = append(q.query.Where, filters...)
	return q
}

// Order sorts the returned records based on the given conditions.
// If multiple conditions are given, they are applied one after the other.
// Note: If OrderRandom is used within the same query,
// it would override the sort conditions.
func (q FieldsLikeDBResponse) Order(by ...*lib.Sort[model.FieldsLikeDBResponse]) FieldsLikeDBResponse {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*lib.SortBuilder)(s))
	}
	return q
}

// OrderRandom sorts the returned records in a random order.
// Note: OrderRandom takes precedence over Order.
func (q FieldsLikeDBResponse) OrderRandom() FieldsLikeDBResponse {
	q.query.SortRandom = true
	return q
}

// Offset skips the first x records for the result set.
func (q FieldsLikeDBResponse) Offset(offset int) FieldsLikeDBResponse {
	q.query.Offset = offset
	return q
}

// Limit restricts the query to return at most x records.
func (q FieldsLikeDBResponse) Limit(limit int) FieldsLikeDBResponse {
	q.query.Limit = limit
	return q
}

// Fetch can be used to return related records.
// This works for both records links and edges.
func (q FieldsLikeDBResponse) Fetch(fetch ...with.Fetch_[model.FieldsLikeDBResponse]) FieldsLikeDBResponse {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
	return q
}

// Timeout adds an execution time limit to the query.
// When exceeded, the query call will return with an error.
func (q FieldsLikeDBResponse) Timeout(timeout time.Duration) FieldsLikeDBResponse {
	q.query.Timeout = timeout
	return q
}

// Parallel tells SurrealDB that individual parts
// of the query can be calculated in parallel.
// This could lead to a faster execution.
func (q FieldsLikeDBResponse) Parallel(parallel bool) FieldsLikeDBResponse {
	q.query.Parallel = parallel
	return q
}

// Count returns the size of the result set, in other words the
// number of records matching the conditions of the query.
func (q FieldsLikeDBResponse) Count(ctx context.Context) (int, error) {
	req := q.query.BuildAsCount()
	raw, err := q.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return 0, err
	}
	var rawCount []queryResult[countResult]
	err = q.unmarshal(raw, &rawCount)
	if err != nil {
		return 0, fmt.Errorf("could not count records: %w", err)
	}
	if len(rawCount) < 1 || len(rawCount[0].Result) < 1 {
		return 0, nil
	}
	return rawCount[0].Result[0].Count, nil
}

// CountAsync is the asynchronous version of Count.
func (q FieldsLikeDBResponse) CountAsync(ctx context.Context) *asyncResult[int] {
	return async(ctx, q.Count)
}

// Exists returns whether at least one record for the conditions
// of the query exists or not. In other words it returns whether
// the size of the result set is greater than 0.
func (q FieldsLikeDBResponse) Exists(ctx context.Context) (bool, error) {
	count, err := q.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsAsync is the asynchronous version of Exists.
func (q FieldsLikeDBResponse) ExistsAsync(ctx context.Context) *asyncResult[bool] {
	return async(ctx, q.Exists)
}

// All returns all records matching the conditions of the query.
func (q FieldsLikeDBResponse) All(ctx context.Context) ([]*model.FieldsLikeDBResponse, error) {
	req := q.query.BuildAsAll()
	res, err := q.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var rawNodes []queryResult[conv.FieldsLikeDBResponse]
	err = q.unmarshal(res, &rawNodes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal records: %w", err)
	}
	if len(rawNodes) < 1 {
		return nil, errors.New("empty result")
	}
	var nodes []*model.FieldsLikeDBResponse
	for _, rawNode := range rawNodes[0].Result {
		node := conv.ToFieldsLikeDBResponse(rawNode)
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

// AllAsync is the asynchronous version of All.
func (q FieldsLikeDBResponse) AllAsync(ctx context.Context) *asyncResult[[]*model.FieldsLikeDBResponse] {
	return async(ctx, q.All)
}

// AllIDs returns the IDs of all records matching the conditions of the query.
func (q FieldsLikeDBResponse) AllIDs(ctx context.Context) ([]string, error) {
	req := q.query.BuildAsAllIDs()
	res, err := q.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var rawNodes []idNode
	err = q.unmarshal(res, &rawNodes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal records: %w", err)
	}
	var ids []string
	for _, rawNode := range rawNodes {
		ids = append(ids, rawNode.ID)
	}
	return ids, nil
}

// AllIDsAsync is the asynchronous version of AllIDs.
func (q FieldsLikeDBResponse) AllIDsAsync(ctx context.Context) *asyncResult[[]string] {
	return async(ctx, q.AllIDs)
}

// First returns the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q FieldsLikeDBResponse) First(ctx context.Context) (*model.FieldsLikeDBResponse, error) {
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
func (q FieldsLikeDBResponse) FirstAsync(ctx context.Context) *asyncResult[*model.FieldsLikeDBResponse] {
	return async(ctx, q.First)
}

// FirstID returns the ID of the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q FieldsLikeDBResponse) FirstID(ctx context.Context) (string, error) {
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
func (q FieldsLikeDBResponse) FirstIDAsync(ctx context.Context) *asyncResult[string] {
	return async(ctx, q.FirstID)
}

// Describe returns a string representation of the query.
// While this might be a valid SurrealDB query, it
// should only be used for debugging purposes.
func (q FieldsLikeDBResponse) Describe() string {
	req := q.query.BuildAsAll()
	return strings.TrimSpace(req.Statement)
}
