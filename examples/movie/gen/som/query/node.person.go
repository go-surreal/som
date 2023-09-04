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
	"strings"
	"time"
)

type nodePerson struct {
	db        Database
	query     lib.Query[model.Person]
	unmarshal func(buf []byte, val any) error
}

type NodePerson struct {
	nodePerson
}

type NodePersonNoLive struct {
	nodePerson
}

func NewPerson(db Database, unmarshal func(buf []byte, val any) error) NodePerson {
	return NodePerson{nodePerson{
		db:        db,
		query:     lib.NewQuery[model.Person]("person"),
		unmarshal: unmarshal,
	}}
}

// Filter adds a where statement to the query to
// select records based on the given conditions.
//
// Use where.All to chain multiple conditions
// together that all need to match.
// Use where.Any to chain multiple conditions
// together where at least one needs to match.
func (q nodePerson) Filter(filters ...lib.Filter[model.Person]) NodePerson {
	q.query.Where = append(q.query.Where, filters...)
	return NodePerson{q}
}

// Order sorts the returned records based on the given conditions.
// If multiple conditions are given, they are applied one after the other.
// Note: If OrderRandom is used within the same query,
// it would override the sort conditions.
func (q nodePerson) Order(by ...*lib.Sort[model.Person]) NodePersonNoLive {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*lib.SortBuilder)(s))
	}
	return NodePersonNoLive{q}
}

// OrderRandom sorts the returned records in a random order.
// Note: OrderRandom takes precedence over Order.
func (q nodePerson) OrderRandom() NodePersonNoLive {
	q.query.SortRandom = true
	return NodePersonNoLive{q}
}

// Offset skips the first x records for the result set.
func (q nodePerson) Offset(offset int) NodePersonNoLive {
	q.query.Offset = offset
	return NodePersonNoLive{q}
}

// Limit restricts the query to return at most x records.
func (q nodePerson) Limit(limit int) NodePersonNoLive {
	q.query.Limit = limit
	return NodePersonNoLive{q}
}

// Fetch can be used to return related records.
// This works for both records links and edges.
func (q nodePerson) Fetch(fetch ...with.Fetch_[model.Person]) NodePerson {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
	return NodePerson{q}
}

// Timeout adds an execution time limit to the query.
// When exceeded, the query call will return with an error.
func (q nodePerson) Timeout(timeout time.Duration) NodePersonNoLive {
	q.query.Timeout = timeout
	return NodePersonNoLive{q}
}

// Parallel tells SurrealDB that individual parts
// of the query can be calculated in parallel.
// This could lead to a faster execution.
func (q nodePerson) Parallel(parallel bool) NodePersonNoLive {
	q.query.Parallel = parallel
	return NodePersonNoLive{q}
}

// Count returns the size of the result set, in other words the
// number of records matching the conditions of the query.
func (q nodePerson) Count(ctx context.Context) (int, error) {
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
func (q nodePerson) CountAsync(ctx context.Context) *asyncResult[int] {
	return async(ctx, q.Count)
}

// Exists returns whether at least one record for the conditions
// of the query exists or not. In other words it returns whether
// the size of the result set is greater than 0.
func (q nodePerson) Exists(ctx context.Context) (bool, error) {
	count, err := q.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsAsync is the asynchronous version of Exists.
func (q nodePerson) ExistsAsync(ctx context.Context) *asyncResult[bool] {
	return async(ctx, q.Exists)
}

// All returns all records matching the conditions of the query.
func (q nodePerson) All(ctx context.Context) ([]*model.Person, error) {
	req := q.query.BuildAsAll()
	res, err := q.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var rawNodes []queryResult[*conv.Person]
	err = q.unmarshal(res, &rawNodes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal records: %w", err)
	}
	if len(rawNodes) < 1 {
		return nil, errors.New("empty result")
	}
	var nodes []*model.Person
	for _, rawNode := range rawNodes[0].Result {
		node := conv.ToPerson(rawNode)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// AllAsync is the asynchronous version of All.
func (q nodePerson) AllAsync(ctx context.Context) *asyncResult[[]*model.Person] {
	return async(ctx, q.All)
}

// AllIDs returns the IDs of all records matching the conditions of the query.
func (q nodePerson) AllIDs(ctx context.Context) ([]string, error) {
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
func (q nodePerson) AllIDsAsync(ctx context.Context) *asyncResult[[]string] {
	return async(ctx, q.AllIDs)
}

// First returns the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q nodePerson) First(ctx context.Context) (*model.Person, error) {
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
func (q nodePerson) FirstAsync(ctx context.Context) *asyncResult[*model.Person] {
	return async(ctx, q.First)
}

// FirstID returns the ID of the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q nodePerson) FirstID(ctx context.Context) (string, error) {
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
func (q nodePerson) FirstIDAsync(ctx context.Context) *asyncResult[string] {
	return async(ctx, q.FirstID)
}

// Live registers the constructed query as a live query.
// Whenever something in the database changes that matches the
// query conditions, the result channel will receive an update.
// If the context is canceled, the result channel will be closed.
func (q NodePerson) Live(ctx context.Context) (<-chan LiveResult[*model.Person], error) {
	req := q.query.BuildAsLive()
	resChan, err := q.db.Live(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query live records: %w", err)
	}
	return live[*model.Person](ctx, resChan, q.unmarshal), nil
}

// Describe returns a string representation of the query.
// While this might be a valid SurrealDB query, it
// should only be used for debugging purposes.
func (q nodePerson) Describe() string {
	req := q.query.BuildAsAll()
	return strings.TrimSpace(req.Statement)
}
