package conv

import model "github.com/marcbinz/som/example/model"

type MemberOfMeta struct {
	IsAdmin  bool `json:"is_admin,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
}

func FromMemberOfMeta(data *model.MemberOfMeta) *MemberOfMeta {
	if data == nil {
		return &MemberOfMeta{}
	}
	return &MemberOfMeta{
		IsActive: data.IsActive,
		IsAdmin:  data.IsAdmin,
	}
}
func ToMemberOfMeta(data *MemberOfMeta) *model.MemberOfMeta {
	return &model.MemberOfMeta{
		IsActive: data.IsActive,
		IsAdmin:  data.IsAdmin,
	}
}
