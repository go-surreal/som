package sdb

import (
	"context"
	"errors"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
)

var User user

type user struct{}

func (user) Query() *query.User {
	return query.NewUser()
}
func (user) Create(ctx context.Context, db *Client, user *model.User) error {
	if user.ID != "" {
		return errors.New("ID must not be set for a node to be created")
	}
	data := conv.FromUser(*user)
	raw, err := db.Create("user", data)
	if err != nil {
		return err
	}
	res := conv.ToUser(raw.([]any)[0].(map[string]any))
	*user = res
	return nil
}
func (user) Read(ctx context.Context, db *Client, id string) (*model.User, error) {
	raw, err := db.db.Select("user:" + id)
	if err != nil {
		return nil, err
	}
	res := conv.ToUser(raw.(map[string]any))
	return &res, nil
}
func (user) Update(ctx context.Context, user *model.User) error {
	return nil
}
func (user) Delete(ctx context.Context, user *model.User) error {
	return nil
}
