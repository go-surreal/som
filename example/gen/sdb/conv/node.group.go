package conv

import (
	"encoding/json"
	model "github.com/marcbinz/som/example/model"
	"strings"
)

type Group struct {
	ID      string     `json:"id,omitempty"`
	Name    string     `json:"name,omitempty"`
	Members []MemberOf `json:"members,omitempty"`
}

func FromGroup(data *model.Group) *Group {
	if data == nil {
		return &Group{}
	}
	return &Group{
		ID:   buildDatabaseID("group", data.ID),
		Name: data.Name,
	}
}
func ToGroup(data *Group) *model.Group {
	return &model.Group{
		ID:   parseDatabaseID("group", data.ID),
		Name: data.Name,
	}
}

type GroupField Group

func (f *GroupField) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}
func (f *GroupField) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\""); strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("group", raw)
		return nil
	}
	type fieldAlias GroupField
	var field fieldAlias
	err := json.Unmarshal(data, &field)
	if err == nil {
		*f = GroupField(field)
	}
	return err
}
func fromGroupField(field GroupField) *model.Group {
	node := Group(field)
	return ToGroup(&node)
}
func toGroupField(node *model.Group) GroupField {
	return GroupField(*FromGroup(node))
}
