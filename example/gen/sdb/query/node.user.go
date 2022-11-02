package query

import (
	"fmt"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	model "github.com/marcbinz/sdb/example/model"
	builder "github.com/marcbinz/sdb/lib/builder"
	filter "github.com/marcbinz/sdb/lib/filter"
	sort "github.com/marcbinz/sdb/lib/sort"
	"time"
)

type User struct {
	db    Database
	query *builder.Query
}

func NewUser(db Database) *User {
	return &User{
		db:    db,
		query: builder.NewQuery("user"),
	}
}
func (q *User) Filter(filters ...filter.Of[model.User]) *User {
	for _, f := range filters {
		q.query.Where = append(q.query.Where, builder.Where(f))
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
	res := builder.Build(q.query)
	fmt.Println(res.Statement, res.Variables)
	raw, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	asMap := raw.([]any)[0].(map[string]any)
	rows := asMap["result"].([]any)
	var nodes []*model.User
	for _, row := range rows {
		node := conv.ToUser(row.(map[string]any))
		nodes = append(nodes, &node)
	}
	return nodes, nil
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
