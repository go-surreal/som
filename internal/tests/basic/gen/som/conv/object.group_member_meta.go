// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import model "github.com/go-surreal/som/tests/basic/model"

type groupMemberMeta struct {
	IsAdmin  bool `cbor:"is_admin"`
	IsActive bool `cbor:"is_active"`
}

func fromGroupMemberMeta(data model.GroupMemberMeta) groupMemberMeta {
	return groupMemberMeta{
		IsActive: data.IsActive,
		IsAdmin:  data.IsAdmin,
	}
}
func fromGroupMemberMetaPtr(data *model.GroupMemberMeta) *groupMemberMeta {
	if data == nil {
		return nil
	}
	return &groupMemberMeta{
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
func toGroupMemberMetaPtr(data *groupMemberMeta) *model.GroupMemberMeta {
	if data == nil {
		return nil
	}
	return &model.GroupMemberMeta{
		IsActive: data.IsActive,
		IsAdmin:  data.IsAdmin,
	}
}
