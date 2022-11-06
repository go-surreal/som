package query

import (
	"errors"
	"fmt"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	with "github.com/marcbinz/sdb/example/gen/sdb/with"
	model "github.com/marcbinz/sdb/example/model"
	builder "github.com/marcbinz/sdb/lib/builder"
	filter "github.com/marcbinz/sdb/lib/filter"
	sort "github.com/marcbinz/sdb/lib/sort"
	"strings"
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
func (q *Group) Fetch(fetch ...with.Fetch_[model.Group]) *Group {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
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
func (q *Group) Count() (int, error) {
	res := q.query.BuildAsCount()
	rows, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return 0, err
	}
	if len(rows) < 1 {
		return 0, nil
	}
	return int(rows[0]["count"].(float64)), nil
}
func (q *Group) Exists() (bool, error) {
	count, err := q.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (q *Group) All() ([]*model.Group, error) {
	res := q.query.BuildAsAll()
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
	res := q.query.BuildAsAllIDs()
	rows, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, row := range rows {
		id := strings.TrimPrefix(row["id"].(string), "group:")
		ids = append(ids, id)
	}
	return ids, nil
}
func (q *Group) First() (*model.Group, error) {
	q.query.Limit = 1
	res, err := q.All()
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return nil, errors.New("empty result")
	}
	return res[0], nil
}
func (q *Group) FirstID() (string, error) {
	q.query.Limit = 1
	res, err := q.AllIDs()
	if err != nil {
		return "", err
	}
	if len(res) < 1 {
		return "", errors.New("empty result")
	}
	return res[0], nil
}
func (q *Group) Describe() (string, error) {
	return "", nil
}
