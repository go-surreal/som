package sdb

import (
	"context"
	"errors"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
)

func (c *Client) User() *user {
	return &user{client: c}
}

type user struct {
	client *Client
}

func (n *user) Query() *query.User {
	return query.NewUser(n.client.db)
}
func (n *user) Create(ctx context.Context, user *model.User) error {
	if user.ID != "" {
		return errors.New("ID must not be set for a node to be created")
	}
	data := conv.FromUser(*user)
	raw, err := n.client.db.Create("user", data)
	if err != nil {
		return err
	}
	res := conv.ToUser(raw.([]any)[0].(map[string]any))
	*user = res
	return nil
}
func (n *user) Read(ctx context.Context, id string) (*model.User, error) {
	raw, err := n.client.db.Select("user" + id)
	if err != nil {
		return nil, err
	}
	res := conv.ToUser(raw.(map[string]any))
	return &res, nil
}
func (n *user) Update(ctx context.Context, user *model.User) error {
	return nil
}
func (n *user) Delete(ctx context.Context, user *model.User) error {
	return nil
}
