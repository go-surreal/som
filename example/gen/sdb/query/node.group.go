package query

import (
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	model "github.com/marcbinz/sdb/example/model"
	builder "github.com/marcbinz/sdb/lib/builder"
	filter "github.com/marcbinz/sdb/lib/filter"
	sort "github.com/marcbinz/sdb/lib/sort"
	"time"
)

type Group struct {
	db    Database
	query *builder.Query
}

func NewGroup(db Database) *Group {
	return &Group{
		db:    db,
		query: builder.NewQuery("group"),
	}
}
func (q *Group) Filter(filters ...filter.Of[model.Group]) *Group {
	for _, f := range filters {
		q.query.Where = append(q.query.Where, builder.Where(f))
	}
	return q
}
func (q *Group) Order(by ...*sort.Of[model.Group]) *Group {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*builder.Sort)(s))
	}
	return q
}
func (q *Group) OrderRandom() *Group {
	q.query.SortRandom = true
	return q
}
func (q *Group) Offset(offset int) *Group {
	q.query.Offset = offset
	return q
}
func (q *Group) Limit(limit int) *Group {
	q.query.Limit = limit
	return q
}
func (q *Group) Unique() *Group {
	return q
}
func (q *Group) Fetch() *Group {
	return q
}
func (q *Group) FetchDepth() *Group {
	return q
}
func (q *Group) Timeout(timeout time.Duration) *Group {
	q.query.Timeout = timeout
	return q
}
func (q *Group) Parallel(parallel bool) *Group {
	q.query.Parallel = parallel
	return q
}
func (q *Group) Count() *Group {
	return q
}
func (q *Group) Exist() *Group {
	return q
}
func (q *Group) All() ([]*model.Group, error) {
	res := builder.Build(q.query)
	rows, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	var nodes []*model.Group
	for _, row := range rows {
		node := conv.ToGroup(row)
		nodes = append(nodes, &node)
	}
	return nodes, nil
}
func (q *Group) AllIDs() ([]string, error) {
	return nil, nil
}
func (q *Group) First() (*model.Group, error) {
	return nil, nil
}
func (q *Group) FirstID() (string, error) {
	return "", nil
}
func (q *Group) Only() (*model.Group, error) {
	return nil, nil
}
func (q *Group) OnlyID() (string, error) {
	return "", nil
}
func (q *Group) Describe() (string, error) {
	return "", nil
}
