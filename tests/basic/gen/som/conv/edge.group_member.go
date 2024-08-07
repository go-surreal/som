// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	sdbc "github.com/go-surreal/sdbc"
	som "github.com/go-surreal/som"
	model "github.com/go-surreal/som/tests/basic/model"
)

type GroupMember struct {
	ID        *sdbc.ID        `cbor:"id,omitempty"`
	CreatedAt *sdbc.DateTime  `cbor:"created_at,omitempty"`
	UpdatedAt *sdbc.DateTime  `cbor:"updated_at,omitempty"`
	Meta      groupMemberMeta `cbor:"meta"`
}

func FromGroupMember(data *model.GroupMember) *GroupMember {
	if data == nil {
		return nil
	}
	return &GroupMember{Meta: noPtrFunc(fromGroupMemberMeta)(data.Meta)}
}

func ToGroupMember(data *GroupMember) *model.GroupMember {
	if data == nil {
		return nil
	}
	return &model.GroupMember{
		Edge:       som.NewEdge(data.ID),
		Meta:       noPtrFunc(toGroupMemberMeta)(data.Meta),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}
