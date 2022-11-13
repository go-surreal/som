package sdb

import (
	"context"
	"errors"
	conv "github.com/marcbinz/sdb/example/gen/sdb/conv"
	query "github.com/marcbinz/sdb/example/gen/sdb/query"
	model "github.com/marcbinz/sdb/example/model"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
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
	data := conv.FromGroup(group)
	raw, err := n.client.db.Create("group", data)
	if err != nil {
		return err
	}
	var convNode conv.Group
	err = surrealdbgo.Unmarshal(raw, &convNode)
	if err != nil {
		return err
	}
	*group = *conv.ToGroup(&convNode)
	return nil
}
func (n *group) Update(ctx context.Context, group *model.Group) error {
	return nil
}
func (n *group) Delete(ctx context.Context, group *model.Group) error {
	return nil
}
