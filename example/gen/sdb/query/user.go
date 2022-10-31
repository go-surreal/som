package query

import (
	model "github.com/marcbinz/sdb/example/model"
	builder "github.com/marcbinz/sdb/lib/builder"
	filter "github.com/marcbinz/sdb/lib/filter"
	sort "github.com/marcbinz/sdb/lib/sort"
	"time"
)

type User struct {
	build *builder.Query
}

func NewUser() *User {
	return &User{build: builder.NewQuery()}
}
func (q *User) Filter(filters ...filter.Of[model.User]) *User {
	for _, f := range filters {
		q.build.Where = append(q.build.Where, builder.Where(f))
	}
	return q
}
func (q *User) Sort(by ...*sort.Of[model.User]) *User {
	return q
}
func (q *User) Offset(offset int) *User {
	return q
}
func (q *User) Limit(limit int) *User {
	return q
}
func (q *User) Unique() *User {
	return q
}
func (q *User) Timeout(timeout time.Duration) *User {
	return q
}
func (q *User) Parallel(parallel bool) *User {
	return q
}
func (q *User) Count() *User {
	return q
}
func (q *User) Exist() *User {
	return q
}
func (q *User) All() ([]*model.User, error) {
	return nil, nil
}
func (q *User) AllIDs() ([]string, error) {
	return nil, nil
}
func (q *User) First() (*model.User, error) {
	return nil, nil
}
func (q *User) FirstID() (string, error) {
	return "", nil
}
func (q *User) Only() (*model.User, error) {
	return nil, nil
}
func (q *User) OnlyID() (string, error) {
	return "", nil
}
