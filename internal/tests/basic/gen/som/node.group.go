// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	som "github.com/go-surreal/som"
	conv "github.com/go-surreal/som/tests/basic/gen/som/conv"
	query "github.com/go-surreal/som/tests/basic/gen/som/query"
	relate "github.com/go-surreal/som/tests/basic/gen/som/relate"
	model "github.com/go-surreal/som/tests/basic/model"
)

type GroupRepo interface {
	Query() query.Builder[model.Group, conv.Group]
	Create(ctx context.Context, group *model.Group) error
	CreateWithID(ctx context.Context, id string, group *model.Group) error
	Read(ctx context.Context, id *som.ID) (*model.Group, bool, error)
	Update(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, group *model.Group) error
	Refresh(ctx context.Context, group *model.Group) error
	Relate() *relate.Group
}

// GroupRepo returns a new repository instance for the Group model.
func (c *ClientImpl) GroupRepo() GroupRepo {
	return &group{repo: &repo[model.Group, conv.Group]{
		db:       c.db,
		name:     "group",
		convTo:   conv.ToGroup,
		convFrom: conv.FromGroup}}
}

type group struct {
	*repo[model.Group, conv.Group]
}

// Query returns a new query builder for the Group model.
func (r *group) Query() query.Builder[model.Group, conv.Group] {
	return query.NewGroup(r.db)
}

// Create creates a new record for the Group model.
// The ID will be generated automatically as a ULID.
func (r *group) Create(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != nil {
		return errors.New("given node already has an id")
	}
	return r.create(ctx, group)
}

// CreateWithID creates a new record for the Group model with the given id.
func (r *group) CreateWithID(ctx context.Context, id string, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != nil {
		return errors.New("given node already has an id")
	}
	return r.createWithID(ctx, id, group)
}

// Read returns the record for the given id, if it exists.
// The returned bool indicates whether the record was found or not.
func (r *group) Read(ctx context.Context, id *som.ID) (*model.Group, bool, error) {
	return r.read(ctx, id)
}

// Update updates the record for the given model.
func (r *group) Update(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() == nil {
		return errors.New("cannot update Group without existing record ID")
	}
	return r.update(ctx, group.ID(), group)
}

// Delete deletes the record for the given model.
func (r *group) Delete(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	return r.delete(ctx, group.ID(), group)
}

// Refresh refreshes the given model with the remote data.
func (r *group) Refresh(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() == nil {
		return errors.New("cannot refresh Group without existing record ID")
	}
	return r.refresh(ctx, group.ID(), group)
}

// Relate returns a new relate instance for the Group model.
func (r *group) Relate() *relate.Group {
	return relate.NewGroup(r.db)
}