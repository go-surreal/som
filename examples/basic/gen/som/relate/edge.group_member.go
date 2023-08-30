// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package relate

import (
	"context"
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/basic/gen/som/conv"
	model "github.com/marcbinz/som/examples/basic/model"
)

type groupMember struct {
	db        Database
	unmarshal func(buf []byte, val any) error
}

func (e groupMember) Create(ctx context.Context, edge *model.GroupMember) error {
	if edge == nil {
		return errors.New("the given edge must not be nil")
	}
	if edge.ID() != "" {
		return errors.New("ID must not be set for an edge to be created")
	}
	if edge.User.ID() == "" {
		return errors.New("ID of the incoming node 'User' must not be empty")
	}
	if edge.Group.ID() == "" {
		return errors.New("ID of the outgoing node 'Group' must not be empty")
	}
	query := "RELATE " + "user:" + edge.User.ID() + "->group_member->" + "group:" + edge.Group.ID() + " CONTENT $data"
	data := conv.FromGroupMember(*edge)
	res, err := e.db.Query(ctx, query, map[string]any{"data": data})
	if err != nil {
		return fmt.Errorf("could not create relation: %w", err)
	}
	var convEdge conv.GroupMember
	err = e.unmarshal(res, &convEdge)
	if err != nil {
		return fmt.Errorf("could not unmarshal relation: %w", err)
	}
	*edge = conv.ToGroupMember(convEdge)
	return nil
}

func (groupMember) Update(edge *model.GroupMember) error {
	return errors.New("not yet implemented")
}

func (groupMember) Delete(edge *model.GroupMember) error {
	return errors.New("not yet implemented")
}
