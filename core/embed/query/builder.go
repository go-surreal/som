///go:build embed

package query

import (
	"context"
	"errors"
	"fmt"
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	with "github.com/go-surreal/som/tests/basic/gen/som/with"
	"io"
	"strings"
	"time"
)

// M is a placeholder for the model type.
// C is a placeholder for the conversion type.
type builder[M, C any] struct {
	db        Database
	query     lib.Query[M]
	unmarshal func(buf []byte, val any) error

	convFrom func(*M) *C
	convTo   func(*C) *M
}

type Builder[M, C any] struct {
	builder[M, C]
}

type BuilderNoLive[M, C any] struct {
	builder[M, C]
}

// Filter adds a where statement to the query to
// select records based on the given conditions.
//
// Use where.All to chain multiple conditions
// together that all need to match.
// Use where.Any to chain multiple conditions
// together where at least one needs to match.
func (b builder[M, C]) Filter(filters ...lib.Filter[M]) Builder[M, C] {
	b.query.Where = append(b.query.Where, filters...)
	return Builder[M, C]{b}
}

// Order sorts the returned records based on the given conditions.
// If multiple conditions are given, they are applied one after the other.
// Note: If OrderRandom is used within the same query,
// it would override the sort conditions.
func (b builder[M, C]) Order(by ...*lib.Sort[M]) BuilderNoLive[M, C] {
	for _, s := range by {
		b.query.Sort = append(b.query.Sort, (*lib.SortBuilder)(s))
	}
	return BuilderNoLive[M, C]{b}
}

// OrderRandom sorts the returned records in a random order.
// Note: OrderRandom takes precedence over Order.
func (b builder[M, C]) OrderRandom() BuilderNoLive[M, C] {
	b.query.SortRandom = true
	return BuilderNoLive[M, C]{b}
}

// Offset skips the first x records for the result set.
func (b builder[M, C]) Offset(offset int) BuilderNoLive[M, C] {
	b.query.Offset = offset
	return BuilderNoLive[M, C]{b}
}

// Limit restricts the query to return at most x records.
func (b builder[M, C]) Limit(limit int) BuilderNoLive[M, C] {
	b.query.Limit = limit
	return BuilderNoLive[M, C]{b}
}

// Fetch can be used to return related records.
// This works for both record links and edges.
func (b builder[M, C]) Fetch(fetch ...with.Fetch_[M]) Builder[M, C] {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			b.query.Fetch = append(b.query.Fetch, field)
		}
	}
	return Builder[M, C]{b}
}

// Timeout adds an execution time limit to the query.
// When exceeded, the query call will return with an error.
func (b builder[M, C]) Timeout(timeout time.Duration) BuilderNoLive[M, C] {
	b.query.Timeout = timeout
	return BuilderNoLive[M, C]{b}
}

// Parallel tells SurrealDB that individual parts
// of the query can be calculated in parallel.
// This could lead to a faster execution.
func (b builder[M, C]) Parallel(parallel bool) BuilderNoLive[M, C] {
	b.query.Parallel = parallel
	return BuilderNoLive[M, C]{b}
}

// Count returns the size of the result set, in other words, the
// number of records matching the conditions of the query.
func (b builder[M, C]) Count(ctx context.Context) (int, error) {
	req := b.query.BuildAsCount()
	raw, err := b.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return 0, err
	}
	var rawCount []queryResult[countResult]
	err = b.unmarshal(raw, &rawCount)
	if err != nil {
		return 0, fmt.Errorf("could not count records: %w", err)
	}
	if len(rawCount) < 1 || len(rawCount[0].Result) < 1 {
		return 0, nil
	}
	return rawCount[0].Result[0].Count, nil
}

// CountAsync is the asynchronous version of Count.
func (b builder[M, C]) CountAsync(ctx context.Context) *asyncResult[int] {
	return async(ctx, b.Count)
}

// Exists returns whether at least one record for the conditions
// of the query exists or not. In other words, it returns whether
// the size of the result set is greater than 0.
func (b builder[M, C]) Exists(ctx context.Context) (bool, error) {
	count, err := b.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsAsync is the asynchronous version of Exists.
func (b builder[M, C]) ExistsAsync(ctx context.Context) *asyncResult[bool] {
	return async(ctx, b.Exists)
}

// All returns all records matching the conditions of the query.
func (b builder[M, C]) All(ctx context.Context) ([]*M, error) {
	return b.all(ctx, b.query.BuildAsAll())
}

func (b builder[M, C]) all(ctx context.Context, req *lib.Result) ([]*M, error) {
	res, err := b.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var rawNodes []queryResult[*C]
	err = b.unmarshal(res, &rawNodes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal records: %w", err)
	}
	if len(rawNodes) < 1 {
		return nil, errors.New("empty result")
	}
	var nodes []*M
	for _, rawNode := range rawNodes[0].Result {
		node := b.convTo(rawNode)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// AllAsync is the asynchronous version of All.
func (b builder[M, C]) AllAsync(ctx context.Context) *asyncResult[[]*M] {
	return async(ctx, b.All)
}

// AllIDs returns the IDs of all records matching the conditions of the query.
func (b builder[M, C]) AllIDs(ctx context.Context) ([]string, error) {
	req := b.query.BuildAsAllIDs()
	res, err := b.db.Query(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query records: %w", err)
	}
	var rawNodes []idNode
	err = b.unmarshal(res, &rawNodes)
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
func (b builder[M, C]) AllIDsAsync(ctx context.Context) *asyncResult[[]string] {
	return async(ctx, b.AllIDs)
}

// First returns the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (b builder[M, C]) First(ctx context.Context) (*M, error) {
	b.query.Limit = 1
	res, err := b.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("empty result")
	}
	return res[0], nil
}

// FirstAsync is the asynchronous version of First.
func (b builder[M, C]) FirstAsync(ctx context.Context) *asyncResult[*M] {
	return async(ctx, b.First)
}

// FirstID returns the ID of the first record matching the conditions of the query.
// This comes in handy when using a filter for a field with unique values or when
// sorting the result set in a specific order where only the first result is relevant.
func (b builder[M, C]) FirstID(ctx context.Context) (string, error) {
	b.query.Limit = 1
	res, err := b.AllIDs(ctx)
	if err != nil {
		return "", err
	}
	if len(res) < 1 {
		return "", errors.New("empty result")
	}
	return res[0], nil
}

// FirstIDAsync is the asynchronous version of FirstID.
func (b builder[M, C]) FirstIDAsync(ctx context.Context) *asyncResult[string] {
	return async(ctx, b.FirstID)
}

// Live registers the constructed query as a live query.
// Whenever something in the database changes that matches the
// query conditions, the result channel will receive an update.
// If the context is canceled, the result channel will be closed.
//
// Note: If you want both the current result set and live updates,
// it is advised to execute the live query first. This is to ensure
// data consistency. The other way around, there could be missing
// updates happening between the initial query and the live query.
func (b builder[M, C]) Live(ctx context.Context) (<-chan LiveResult[*M], error) {
	req := b.query.BuildAsLive()
	resChan, err := b.db.Live(ctx, req.Statement, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("could not query live records: %w", err)
	}
	return live(ctx, resChan, b.unmarshal, b.convTo), nil
}

// LiveDiff behaves like Live, but instead of receiving the full result
// set on every change, it only receives the actual changes.
//func (b builder[M, C]) LiveDiff(ctx context.Context) (<-chan LiveResult[*M], error) {
//	panic("not yet implemented") // TODO: implement!
//}

// Describe returns a string representation of the query.
// While this might be a valid SurrealDB query, it
// should only be used for debugging purposes.
func (b builder[M, C]) Describe() string {
	req := b.query.BuildAsAll()
	return strings.TrimSpace(req.Statement)
}

func (b builder[M, C]) Iterator(batch int) *Iterator[M] {
	return &Iterator[M]{
		partialQuery: b.query,
		executor:     b.all,

		batch:  batch,
		offset: b.query.Offset,
		limit:  b.query.Limit,
	}
}

type Iterator[M any] struct {
	partialQuery lib.Query[M]
	executor     func(ctx context.Context, req *lib.Result) ([]*M, error)

	current []*M
	index   int

	batch  int
	offset int
	limit  int
}

func (i *Iterator[M]) Next() (*M, error) {
	if i.index < len(i.current) {
		val := i.current[i.index]
		i.index++
		return val, nil
	}

	query := i.partialQuery
	query.Offset = i.offset
	query.Limit = i.limit

	var err error
	i.current, err = query.BuildAsAll()
	if err != nil {
		var t T
		return t, err
	}

	if len(i.current) == 0 {
		var t T
		return t, io.EOF
	}

	i.index = 0
	return i.Next()
}
