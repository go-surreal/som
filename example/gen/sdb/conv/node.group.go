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
