package sdb

import (
	"context"
	"errors"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
)

func (c *Client) Group() *group {
	return &group{client: c}
}

type group struct {
	client *Client
}

func (n *group) Query() *query.Group {
	return query.NewGroup(n.client.db)
}
func (n *group) Create(ctx context.Context, group *model.Group) error {
	if group.ID != "" {
		return errors.New("ID must not be set for a node to be created")
	}
	data := conv.FromGroup(*group)
	raw, err := n.client.db.Create("group", data)
	if err != nil {
		return err
	}
	res := conv.ToGroup(raw.([]any)[0].(map[string]any))
	*group = res
	return nil
}
func (n *group) Read(ctx context.Context, id string) (*model.Group, error) {
	raw, err := n.client.db.Select("group" + id)
	if err != nil {
		return nil, err
	}
	res := conv.ToGroup(raw.(map[string]any))
	return &res, nil
}
func (n *group) Update(ctx context.Context, group *model.Group) error {
	return nil
}
func (n *group) Delete(ctx context.Context, group *model.Group) error {
	return nil
}
