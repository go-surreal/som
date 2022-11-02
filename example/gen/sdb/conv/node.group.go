package conv

import model "github.com/marcbinz/sdb/example/model"

func FromGroup(data model.Group) map[string]any {
	return map[string]any{"name": data.Name}
}
func ToGroup(data map[string]any) model.Group {
	return model.Group{
		ID:   prepareID("group", data["id"]),
		Name: data["name"].(string),
	}
}
func fromGroupRecord(data any) model.Group {
	if node, ok := data.(map[string]any); ok {
		return ToGroup(node)
	}
	return model.Group{}
}
func toGroupRecord(node model.Group) any {
	if node.ID == "" {
		return nil
	}
	return "group:" + node.ID
}
