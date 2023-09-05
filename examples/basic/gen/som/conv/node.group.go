// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	"encoding/json"
	som "github.com/go-surreal/som"
	model "github.com/go-surreal/som/examples/basic/model"
	"strings"
	"time"
)

type Group struct {
	ID        string        `json:"id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Name      string        `json:"name"`
	Members   []GroupMember `json:"members,omitempty"`
}

func FromGroup(data *model.Group) *Group {
	if data == nil {
		return nil
	}
	return &Group{Name: data.Name}
}

func ToGroup(data *Group) *model.Group {
	if data == nil {
		return nil
	}
	return &model.Group{
		Members:    mapSlice(data.Members, noPtrFunc(ToGroupMember)),
		Name:       data.Name,
		Node:       som.NewNode(parseDatabaseID("group", data.ID)),
		Timestamps: som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
	}
}

type groupLink struct {
	Group
	ID string
}

func (f *groupLink) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}

func (f *groupLink) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("group", raw)
		return nil
	}
	type alias groupLink
	var link alias
	err := json.Unmarshal(data, &link)
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
	out := ToGroup(&res)
	return *out
}

func fromGroupLinkPtr(link *groupLink) *model.Group {
	if link == nil {
		return nil
	}
	res := Group(link.Group)
	return ToGroup(&res)
}

func toGroupLink(node model.Group) *groupLink {
	if node.ID() == "" {
		return nil
	}
	link := groupLink{Group: *FromGroup(&node), ID: buildDatabaseID("group", node.ID())}
	return &link
}

func toGroupLinkPtr(node *model.Group) *groupLink {
	if node == nil || node.ID() == "" {
		return nil
	}
	link := groupLink{Group: *FromGroup(node), ID: buildDatabaseID("group", node.ID())}
	return &link
}
