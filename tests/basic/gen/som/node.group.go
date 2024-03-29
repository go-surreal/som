// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	conv "github.com/go-surreal/som/tests/basic/gen/som/conv"
	query "github.com/go-surreal/som/tests/basic/gen/som/query"
	relate "github.com/go-surreal/som/tests/basic/gen/som/relate"
	model "github.com/go-surreal/som/tests/basic/model"
)

type GroupRepo interface {
	Query() query.Builder[model.Group, conv.Group]
	Create(ctx context.Context, user *model.Group) error
	CreateWithID(ctx context.Context, id string, user *model.Group) error
	Read(ctx context.Context, id string) (*model.Group, bool, error)
	Update(ctx context.Context, user *model.Group) error
	Delete(ctx context.Context, user *model.Group) error
	Refresh(ctx context.Context, user *model.Group) error
	Relate() *relate.Group
}

func (c *ClientImpl) GroupRepo() GroupRepo {
	return &group{repo: &repo[model.Group, conv.Group]{
		db:        c.db,
		marshal:   c.marshal,
		unmarshal: c.unmarshal,
		name:      "group",
		convTo:    conv.ToGroup,
		convFrom:  conv.FromGroup}}
}

type group struct {
	*repo[model.Group, conv.Group]
}

func (r *group) Query() query.Builder[model.Group, conv.Group] {
	return query.NewGroup(r.db, r.unmarshal)
}

func (r *group) Create(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != "" {
		return errors.New("given node already has an id")
	}
	return r.create(ctx, group)
}

func (r *group) CreateWithID(ctx context.Context, id string, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != "" {
		return errors.New("given node already has an id")
	}
	return r.createWithID(ctx, id, group)
}

func (r *group) Read(ctx context.Context, id string) (*model.Group, bool, error) {
	return r.read(ctx, id)
}

func (r *group) Update(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() == "" {
		return errors.New("cannot update Group without existing record ID")
	}
	return r.update(ctx, group.ID(), group)
}

func (r *group) Delete(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	return r.delete(ctx, group.ID(), group)
}

func (r *group) Refresh(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() == "" {
		return errors.New("cannot refresh Group without existing record ID")
	}
	return r.refresh(ctx, group.ID(), group)
}

func (r *group) Relate() *relate.Group {
	return relate.NewGroup(r.db, r.unmarshal)
}
