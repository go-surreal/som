// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package relate

import (
	"errors"
	"fmt"
	conv "github.com/marcbinz/som/examples/live/gen/som/conv"
	model "github.com/marcbinz/som/examples/live/model"
	surrealdbgo "github.com/surrealdb/surrealdb.go"
)

type groupMember struct {
	db Database
}

func (e groupMember) Create(edge *model.GroupMember) error {
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
	convEdge, err := surrealdbgo.SmartUnmarshal[conv.GroupMember](e.db.Query(query, map[string]any{"data": data}))
	if err != nil {
		return fmt.Errorf("could not create relation: %w", err)
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
