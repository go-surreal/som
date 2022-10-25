package sdb

import (
	"context"
	model "github.com/marcbinz/sdb/db/model"
	query "go.alfnet.dev/service/gampi/gen/sdb/query"
)

var User user

type user struct{}

func (user) Query() *query.User {
	return &query.User{}
}
func (user) Create(ctx context.Context, user *model.User) error {
	return nil
}
func (user) Update(ctx context.Context, user *model.User) error {
	return nil
}
func (user) Delete(ctx context.Context, user *model.User) error {
	return nil
}
