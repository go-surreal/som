// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import (
	som "github.com/marcbinz/som"
	model "github.com/marcbinz/som/example/model"
	"time"
)

type MemberOf struct {
	ID        string       `json:"id,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Meta      memberOfMeta `json:"meta"`
}

func FromMemberOf(data model.MemberOf) MemberOf {
	return MemberOf{Meta: fromMemberOfMeta(data.Meta)}
}
func ToMemberOf(data MemberOf) model.MemberOf {
	return model.MemberOf{
		Edge:       som.NewEdge(parseDatabaseID("member_of", data.ID)),
		Meta:       toMemberOfMeta(data.Meta),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}
