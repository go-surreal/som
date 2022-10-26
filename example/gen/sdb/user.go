package sdb

import (
	"context"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
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
