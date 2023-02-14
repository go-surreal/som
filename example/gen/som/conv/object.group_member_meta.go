// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import model "github.com/marcbinz/som/example/model"

type groupMemberMeta struct {
	IsAdmin  bool `json:"is_admin"`
	IsActive bool `json:"is_active"`
}

func fromGroupMemberMeta(data model.GroupMemberMeta) groupMemberMeta {
	return groupMemberMeta{
		IsActive: data.IsActive,
		IsAdmin:  data.IsAdmin,
	}
}
func toGroupMemberMeta(data groupMemberMeta) model.GroupMemberMeta {
	return model.GroupMemberMeta{
		IsActive: data.IsActive,
		IsAdmin:  data.IsAdmin,
	}
}