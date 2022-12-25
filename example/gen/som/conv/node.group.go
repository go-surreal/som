// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import (
	"encoding/json"
	model "github.com/marcbinz/som/example/model"
	"strings"
)

type Group struct {
	ID      string     `json:"id,omitempty"`
	Name    string     `json:"name"`
	Members []MemberOf `json:"members,omitempty"`
}

func FromGroup(data model.Group) Group {
	return Group{Name: data.Name}
}
func ToGroup(data Group) model.Group {
	return model.Group{
		ID:      parseDatabaseID("group", data.ID),
		Members: mapSlice(data.Members, ToMemberOf),
		Name:    data.Name,
	}
}

type groupLink Group

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
func fromGroupLink(link groupLink) model.Group {
	return ToGroup(Group(link))
}
func toGroupLink(node model.Group) groupLink {
	return groupLink(FromGroup(node))
}
