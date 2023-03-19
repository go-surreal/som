// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package query

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/example/gen/som/conv"
	with "github.com/marcbinz/som/example/gen/som/with"
	model "github.com/marcbinz/som/example/model"
	lib "github.com/marcbinz/som/lib"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
	"strings"
	"time"
)

type Group struct {
	db    Database
	query lib.Query[model.Group]
}

func NewGroup(db Database) Group {
	return Group{
		db:    db,
		query: lib.NewQuery[model.Group]("group"),
	}
}

// Filter adds a where statement to the query to
// select records based on the given conditions.
//
// Use where.All to chain multiple conditions
// together that all need to match.
// Use where.Any to chain multiple conditions
// together where at least one needs to match.
func (q Group) Filter(filters ...lib.Filter[model.Group]) Group {
	q.query.Where = append(q.query.Where, filters...)
	return q
}

// Order sorts the returned records based on the given conditions.
// If multiple conditions are given, they are applied one after the other.
// Note: If OrderRandom is used within the same query,
// it would override the sort conditions.
func (q Group) Order(by ...*lib.Sort[model.Group]) Group {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*lib.SortBuilder)(s))
	}
	return q
}

// OrderRandom sorts the returned records in a random order.
// Note: OrderRandom takes precedence over Order.
func (q Group) OrderRandom() Group {
	q.query.SortRandom = true
	return q
}

// Offset skips the first x records for the result set.
func (q Group) Offset(offset int) Group {
	q.query.Offset = offset
	return q
}

// Limit restricts the query to return at most x records.
func (q Group) Limit(limit int) Group {
	q.query.Limit = limit
	return q
}

// Fetch can be used to return related records.
// This works for both records links and edges.
func (q Group) Fetch(fetch ...with.Fetch_[model.Group]) Group {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
	return q
}

// Timeout adds an execution time limit to the query.
// When exceeded, the query call will return with an error.
func (q Group) Timeout(timeout time.Duration) Group {
	q.query.Timeout = timeout
	return q
}

// Parallel tells SurrealDB that individual parts
// of the query can be calculated in parallel.
// This could lead to a faster execution.
func (q Group) Parallel(parallel bool) Group {
	q.query.Parallel = parallel
	return q
}

// Count returns the size of the result set, in other words the
// number of records matching the conditions of the query.
func (q Group) Count(ctx context.Context) (int, error) {
	res := q.query.BuildAsCount()
	raw, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return 0, err
	}
	var rawCount countResult
	ok, err := surrealdbgo.UnmarshalRaw(raw, &rawCount)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return rawCount.Count, nil
}

// Exists returns whether at least one record for the conditons
// of the query exists or not. In other words it returns whether
// the size of the result set is greater than 0.
func (q Group) Exists(ctx context.Context) (bool, error) {
	count, err := q.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// All returns all records matching the conditions of the query.
func (q Group) All(ctx context.Context) ([]*model.Group, error) {
	res := q.query.BuildAsAll()
	raw, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	var rawNodes []conv.Group
	ok, err := surrealdbgo.UnmarshalRaw(raw, &rawNodes)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	var nodes []*model.Group
	for _, rawNode := range rawNodes {
		node := conv.ToGroup(rawNode)
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

// AllIDs returns the IDs of all records matching the conditions of the query.
func (q Group) AllIDs(ctx context.Context) ([]string, error) {
	res := q.query.BuildAsAllIDs()
	raw, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	var rawNodes []*idNode
	ok, err := surrealdbgo.UnmarshalRaw(raw, &rawNodes)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	var ids []string
	for _, rawNode := range rawNodes {
		ids = append(ids, rawNode.ID)
	}
	return ids, nil
}

// First returns the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q Group) First(ctx context.Context) (*model.Group, error) {
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

// FirstID returns the ID of the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (q Group) FirstID(ctx context.Context) (string, error) {
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

// Describe returns a string representation of the query.
// While this might be a valid SurrealDB query, it
// should only be used for debugging purposes.
func (q Group) Describe() string {
	res := q.query.BuildAsAll()
	return strings.TrimSpace(res.Statement)
}
