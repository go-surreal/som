// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package som

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/basic/gen/som/conv"
	query "github.com/marcbinz/som/examples/basic/gen/som/query"
	relate "github.com/marcbinz/som/examples/basic/gen/som/relate"
	model "github.com/marcbinz/som/examples/basic/model"
	constants "github.com/surrealdb/surrealdb.go/pkg/constants"
	marshal "github.com/surrealdb/surrealdb.go/pkg/marshal"
	"time"
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
	return &group{db: c.db}
}

type group struct {
	db Database
}

func (n *group) Query() query.NodeGroup {
	return query.NewGroup(n.db)
}

func (n *group) Create(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() != "" {
		return errors.New("given node already has an id")
	}
	key := "group"
	data := conv.FromGroup(*group)
	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt
	raw, err := n.db.Create(key, data)
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	var convNodes []conv.Group
	err = marshal.Unmarshal(raw, &convNodes)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}
	if len(convNodes) < 1 {
		return errors.New("response is empty")
	}
	*group = conv.ToGroup(convNodes[0])
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
	data := conv.FromGroup(*group)
	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt
	convNode, err := marshal.SmartUnmarshal[conv.Group](n.db.Create(key, data))
	if err != nil {
		return fmt.Errorf("could not create entity: %w", err)
	}
	*group = conv.ToGroup(convNode[0])
	return nil
}

func (n *group) Read(ctx context.Context, id string) (*model.Group, bool, error) {
	convNode, err := marshal.SmartUnmarshal[conv.Group](n.db.Select("group:⟨" + id + "⟩"))
	if errors.Is(err, constants.ErrNoRow) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("could not read entity: %w", err)
	}
	node := conv.ToGroup(convNode[0])
	return &node, true, nil
}

func (n *group) Update(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	if group.ID() == "" {
		return errors.New("cannot update Group without existing record ID")
	}
	data := conv.FromGroup(*group)
	data.UpdatedAt = time.Now()
	convNode, err := marshal.SmartUnmarshal[conv.Group](n.db.Update("group:⟨"+group.ID()+"⟩", data))
	if err != nil {
		return fmt.Errorf("could not update entity: %w", err)
	}
	*group = conv.ToGroup(convNode[0])
	return nil
}

func (n *group) Delete(ctx context.Context, group *model.Group) error {
	if group == nil {
		return errors.New("the passed node must not be nil")
	}
	_, err := n.db.Delete("group:⟨" + group.ID() + "⟩")
	if err != nil {
		return fmt.Errorf("could not delete entity: %w", err)
	}
	return nil
}

func (n *group) Relate() *relate.Group {
	return relate.NewGroup(n.db)
}
