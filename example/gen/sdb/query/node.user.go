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
func (q *User) Order(by ...*sort.Of[model.User]) *User {
	for _, s := range by {
		q.query.Sort = append(q.query.Sort, (*builder.Sort)(s))
	}
	return q
}
func (q *User) OrderRandom() *User {
	q.query.SortRandom = true
	return q
}
func (q *User) Offset(offset int) *User {
	q.query.Offset = offset
	return q
}
func (q *User) Limit(limit int) *User {
	q.query.Limit = limit
	return q
}
func (q *User) Fetch(fetch ...with.Fetch_[model.User]) *User {
	for _, f := range fetch {
		if field := fmt.Sprintf("%v", f); field != "" {
			q.query.Fetch = append(q.query.Fetch, field)
		}
	}
	return q
}
func (q *User) Timeout(timeout time.Duration) *User {
	q.query.Timeout = timeout
	return q
}
func (q *User) Parallel(parallel bool) *User {
	q.query.Parallel = parallel
	return q
}
func (q *User) Count() (int, error) {
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
func (q *User) Exists() (bool, error) {
	count, err := q.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (q *User) All() ([]*model.User, error) {
	res := q.query.BuildAsAll()
	rows, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	var nodes []*model.User
	for _, row := range rows {
		node := conv.ToUser(row)
		nodes = append(nodes, &node)
	}
	return nodes, nil
}
func (q *User) AllIDs() ([]string, error) {
	res := q.query.BuildAsAllIDs()
	rows, err := q.db.Query(res.Statement, res.Variables)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, row := range rows {
		id := strings.TrimPrefix(row["id"].(string), "user:")
		ids = append(ids, id)
	}
	return ids, nil
}
func (q *User) First() (*model.User, error) {
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
func (q *User) FirstID() (string, error) {
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
func (q *User) Describe() (string, error) {
	return "", nil
}
