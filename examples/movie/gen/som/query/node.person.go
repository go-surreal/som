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
	surrealdbgo "github.com/surrealdb/surrealdb.go"
	"strings"
	"time"
)

type Person struct {
	db    Database
	query lib.Query[model.Person]
}

func NewPerson(db Database) Person {
	return Person{
		db:    db,
		query: lib.NewQuery[model.Person]("person"),
	}
}

// Filter adds a where statement to the query to
// select records based on the given conditions.
//
// Use where.All to chain multiple conditions
// together that all need to match.
// Use where.Any to chain multiple conditions
// together where at least one needs to match.
func (q Person) Filter(filters ...lib.Filter[model.Person]) Person {
	q.query.Where = append(q.query.Where, filters...)
	return q
}

// Order sorts the returned records based on the given conditions.
// If multiple conditions are given, they are applied one after the other.
// Note: If OrderRandom is used within the same query,
// it would override the sort conditions.
func (q Person) Order(by ...*lib.Sort[model.Person]) Person {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*lib.SortBuilder)(s))
	}
	return q
}

// OrderRandom sorts the returned records in a random order.
// Note: OrderRandom takes precedence over Order.
func (q Person) OrderRandom() Person {
	q.query.SortRandom = true
	return q
}

// Offset skips the first x records for the result set.
func (q Person) Offset(offset int) Person {
	q.query.Offset = offset
	return q
}

// Limit restricts the query to return at most x records.
func (q Person) Limit(limit int) Person {
	q.query.Limit = limit
	return q
}

// Fetch can be used to return related records.
// This works for both records links and edges.
func (q Person) Fetch(fetch ...with.Fetch_[model.Person]) Person {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
	return q
}

// Timeout adds an execution time limit to the query.
// When exceeded, the query call will return with an error.
func (q Person) Timeout(timeout time.Duration) Person {
	q.query.Timeout = timeout
	return q
}

// Parallel tells SurrealDB that individual parts
// of the query can be calculated in parallel.
// This could lead to a faster execution.
func (q Person) Parallel(parallel bool) Person {
	q.query.Parallel = parallel
	return q
}

// Count returns the size of the result set, in other words the
// number of records matching the conditions of the query.
func (q Person) Count(ctx context.Context) (int, error) {
	res := q.query.BuildAsCount()
	result, err := surrealdbgo.SmartUnmarshal[[]countResult](q.db.Query(res.Statement, res.Variables))
	if err != nil {
		return 0, fmt.Errorf("could not count records: %w", err)
	}
	if len(result) < 1 {
		return 0, errors.New("database result is empty")
	}
	return result[0].Count, nil
}

// Exists returns whether at least one record for the conditons
// of the query exists or not. In other words it returns whether
// the size of the result set is greater than 0.
func (q Person) Exists(ctx context.Context) (bool, error) {
	count, err := q.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// All returns all records matching the conditions of the query.
func (q Person) All(ctx context.Context) ([]*model.Person, error) {
	res := q.query.BuildAsAll()
	rawNodes, err := surrealdbgo.SmartUnmarshal[[]conv.Person](q.db.Query(res.Statement, res.Variables))
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var nodes []*model.Person
	for _, rawNode := range rawNodes {
		node := conv.ToPerson(rawNode)
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

// AllIDs returns the IDs of all records matching the conditions of the query.
func (q Person) AllIDs(ctx context.Context) ([]string, error) {
	res := q.query.BuildAsAllIDs()
	rawNodes, err := surrealdbgo.SmartUnmarshal[[]idNode](q.db.Query(res.Statement, res.Variables))
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
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
func (q Person) First(ctx context.Context) (*model.Person, error) {
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
func (q Person) FirstID(ctx context.Context) (string, error) {
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
func (q Person) Describe() string {
	res := q.query.BuildAsAll()
	return strings.TrimSpace(res.Statement)
}
