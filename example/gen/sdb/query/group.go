package query

import (
	model "github.com/marcbinz/sdb/example/model"
	builder "github.com/marcbinz/sdb/lib/builder"
	filter "github.com/marcbinz/sdb/lib/filter"
	sort "github.com/marcbinz/sdb/lib/sort"
	"time"
)

type Group struct {
	build *builder.Query
}

func NewGroup() *Group {
	return &Group{build: builder.NewQuery()}
}
func (q *Group) Filter(filters ...filter.Of[model.Group]) *Group {
	for _, f := range filters {
		q.build.Where = append(q.build.Where, builder.Where(f))
	}
	return q
}
func (q *Group) Sort(by ...*sort.Of[model.Group]) *Group {
	return q
}
func (q *Group) Offset(offset int) *Group {
	return q
}
func (q *Group) Limit(limit int) *Group {
	return q
}
func (q *Group) Unique() *Group {
	return q
}
func (q *Group) Timeout(timeout time.Duration) *Group {
	return q
}
func (q *Group) Parallel(parallel bool) *Group {
	return q
}
func (q *Group) Count() *Group {
	return q
}
func (q *Group) Exist() *Group {
	return q
}
func (q *Group) All() ([]*model.Group, error) {
	return nil, nil
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
