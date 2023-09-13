// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/go-surreal/som/examples/basic/gen/som/conv"
	query "github.com/go-surreal/som/examples/basic/gen/som/query"
	relate "github.com/go-surreal/som/examples/basic/gen/som/relate"
	model "github.com/go-surreal/som/examples/basic/model"
)

type GroupRepo interface {
	Query() query.NodeGroup
	Create(ctx context.Context, user *model.Group) error
	CreateWithID(ctx context.Context, id string, user *model.Group) error
	Read(ctx context.Context, id string) (*model.Group, bool, error)
	Update(ctx context.Context, user *model.Group) error
	Delete(ctx context.Context, user *model.Group) error
	Relate() *relate.Group
}

func (c *ClientImpl) GroupRepo() GroupRepo {
	return &group{db: c.db, marshal: c.marshal, unmarshal: c.unmarshal}
}

type group struct {
	db        Database
	marshal   func(val any) ([]byte, error)
	unmarshal func(buf []byte, val any) error
}

func (n *group) Query() query.NodeGroup {
	return query.NewGroup(n.db, n.unmarshal)
}

func (n *group) Create(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != "" {
		return errors.New("given node already has an id")
	}
	key := "group:ulid()"
	data := conv.FromGroup(group)
	raw, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNode *conv.Group
	err = n.unmarshal(raw, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	*group = *conv.ToGroup(convNode)
	return nil
}

func (n *group) CreateWithID(ctx context.Context, id string, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != "" {
		return errors.New("creating node with preset ID not allowed, use CreateWithID for that")
	}
	key := "group:" + "⟨" + id + "⟩"
	data := conv.FromGroup(group)
	res, err := n.db.Create(ctx, key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNode *conv.Group
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*group = *conv.ToGroup(convNode)
	return nil
}

func (n *group) Read(ctx context.Context, id string) (*model.Group, bool, error) {
	res, err := n.db.Select(ctx, "group:⟨"+id+"⟩")
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	var convNode *conv.Group
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return nil, false, fmt.Errorf("could not unmarshal entity: %w", err)
	}
	return conv.ToGroup(convNode), true, nil
}

func (n *group) Update(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() == "" {
		return errors.New("cannot update Group without existing record ID")
	}
	data := conv.FromGroup(group)
	res, err := n.db.Update(ctx, "group:⟨"+group.ID()+"⟩", data)
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	var convNode *conv.Group
	err = n.unmarshal(res, &convNode)
	if err != nil {
		return fmt.Errorf("could not unmarshal entity: %w", err)
	}
	*group = *conv.ToGroup(convNode)
	return nil
}

func (n *group) Delete(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.db.Delete(ctx, "group:⟨"+group.ID()+"⟩")
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}

func (n *group) Relate() *relate.Group {
	return relate.NewGroup(n.db, n.unmarshal)
}