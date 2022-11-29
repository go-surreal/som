package conv

import (
	model "github.com/marcbinz/som/example/model"
	"time"
)

type MemberOf struct {
	ID        string       `json:"id,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
	Meta      MemberOfMeta `json:"meta,omitempty"`
}

func FromMemberOf(data *model.MemberOf) *MemberOf {
	if data == nil {
		return &MemberOf{}
	}
	return &MemberOf{
		CreatedAt: data.CreatedAt,
		ID:        buildDatabaseID("member_of", data.ID),
		Meta:      *FromMemberOfMeta(&data.Meta),
		UpdatedAt: data.UpdatedAt,
	}
}
func ToMemberOf(data *MemberOf) *model.MemberOf {
	return &model.MemberOf{
		CreatedAt: data.CreatedAt,
		ID:        parseDatabaseID("member_of", data.ID),
		Meta:      *ToMemberOfMeta(&data.Meta),
		UpdatedAt: data.UpdatedAt,
	}
}
