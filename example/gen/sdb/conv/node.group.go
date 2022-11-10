package conv

import model "github.com/marcbinz/sdb/example/model"

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FromGroup(data *model.Group) *Group {
	if data == nil {
		return &Group{}
	}
	return &Group{Name: data.Name}
}
func ToGroup(data *Group) *model.Group {
	return &model.Group{
		ID:   prepareID("group", data.ID),
		Name: data.Name,
	}
}
func fromGroupRecord(data any) *model.Group {
	if node, ok := data.(*Group); ok {
		return ToGroup(node)
	}
	return &model.Group{}
}
func toGroupRecord(node model.Group) string {
	if node.ID == "" {
		return ""
	}
	return "group:" + node.ID
}
