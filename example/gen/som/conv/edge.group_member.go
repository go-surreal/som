// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import (
	som "github.com/marcbinz/som"
	model "github.com/marcbinz/som/example/model"
	"time"
)

type GroupMember struct {
	ID        string          `json:"id,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Meta      groupMemberMeta `json:"meta"`
}

func FromGroupMember(data model.GroupMember) GroupMember {
	return GroupMember{Meta: fromGroupMemberMeta(data.Meta)}
}

func ToGroupMember(data GroupMember) model.GroupMember {
	return model.GroupMember{
		Edge:       som.NewEdge(parseDatabaseID("group_member", data.ID)),
		Meta:       toGroupMemberMeta(data.Meta),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}
