// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	sdbc "github.com/go-surreal/sdbc"
	som "github.com/go-surreal/som/tests/basic/gen/som"
	model "github.com/go-surreal/som/tests/basic/model"
)

type GroupMember struct {
	ID        *som.ID         `cbor:"id,omitempty"`
	CreatedAt *sdbc.DateTime  `cbor:"created_at,omitempty"`
	UpdatedAt *sdbc.DateTime  `cbor:"updated_at,omitempty"`
	Meta      groupMemberMeta `cbor:"meta"`
}

func FromGroupMember(data model.GroupMember) GroupMember {
	return GroupMember{Meta: fromGroupMemberMeta(data.Meta)}
}
func FromGroupMemberPtr(data *model.GroupMember) *GroupMember {
	if data == nil {
		return nil
	}
	return &GroupMember{Meta: fromGroupMemberMeta(data.Meta)}
}

func ToGroupMember(data *GroupMember) model.GroupMember {
	return model.GroupMember{
		Edge:       som.NewEdge(data.ID),
		Meta:       toGroupMemberMeta(data.Meta),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}
func ToGroupMemberPtr(data *GroupMember) *model.GroupMember {
	if data == nil {
		return nil
	}
	return &model.GroupMember{
		Edge:       som.NewEdge(data.ID),
		Meta:       toGroupMemberMeta(data.Meta),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}
