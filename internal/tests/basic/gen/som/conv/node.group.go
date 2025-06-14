// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	v2 "github.com/fxamacker/cbor/v2"
	sdbc "github.com/go-surreal/sdbc"
	som "github.com/go-surreal/som/tests/basic/gen/som"
	model "github.com/go-surreal/som/tests/basic/model"
)

type Group struct {
	ID        *som.ID        `cbor:"id,omitempty"`
	CreatedAt *sdbc.DateTime `cbor:"created_at,omitempty"`
	UpdatedAt *sdbc.DateTime `cbor:"updated_at,omitempty"`
	Name      string         `cbor:"name"`
	Members   []*GroupMember `cbor:"members,omitempty"`
}

func FromGroup(data model.Group) Group {
	return Group{Name: data.Name}
}
func FromGroupPtr(data *model.Group) *Group {
	if data == nil {
		return nil
	}
	return &Group{Name: data.Name}
}

func ToGroup(data Group) model.Group {
	return model.Group{
		Members:    mapSliceFn(ToGroupMember)(data.Members),
		Name:       data.Name,
		Node:       som.NewNode(data.ID),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}
func ToGroupPtr(data *Group) *model.Group {
	if data == nil {
		return nil
	}
	return &model.Group{
		Members:    mapSliceFn(ToGroupMember)(data.Members),
		Name:       data.Name,
		Node:       som.NewNode(data.ID),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}

type groupLink struct {
	Group
	ID *som.ID
}

func (f *groupLink) MarshalCBOR() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return v2.Marshal(f.ID)
}

func (f *groupLink) UnmarshalCBOR(data []byte) error {
	if err := v2.Unmarshal(data, &f.ID); err == nil {
		return nil
	}
	type alias groupLink
	var link alias
	err := v2.Unmarshal(data, &link)
	if err == nil {
		*f = groupLink(link)
	}
	return err
}

func fromGroupLink(link *groupLink) model.Group {
	if link == nil {
		return model.Group{}
	}
	res := Group(link.Group)
	return ToGroup(res)
}

func fromGroupLinkPtr(link *groupLink) *model.Group {
	if link == nil {
		return nil
	}
	res := Group(link.Group)
	out := ToGroup(res)
	return &out
}

func toGroupLink(node model.Group) *groupLink {
	if node.ID() == nil {
		return nil
	}
	link := groupLink{Group: FromGroup(node), ID: node.ID()}
	return &link
}

func toGroupLinkPtr(node *model.Group) *groupLink {
	if node == nil || node.ID() == nil {
		return nil
	}
	link := groupLink{Group: FromGroup(*node), ID: node.ID()}
	return &link
}
