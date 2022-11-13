package conv

import (
	model "github.com/marcbinz/sdb/example/model"
	"time"
)

type MemberOf struct {
	Since time.Time    `json:"since,omitempty"`
	Meta  MemberOfMeta `json:"meta,omitempty"`
}

func FromMemberOf(data *model.MemberOf) *MemberOf {
	if data == nil {
		return &MemberOf{}
	}
	return &MemberOf{
		Meta:  *FromMemberOfMeta(&data.Meta),
		Since: data.Since,
	}
}
func ToMemberOf(data *MemberOf) *model.MemberOf {
	return &model.MemberOf{
		Meta:  *ToMemberOfMeta(&data.Meta),
		Since: data.Since,
	}
}
