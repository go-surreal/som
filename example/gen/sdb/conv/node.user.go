package conv

import model "github.com/marcbinz/sdb/example/model"

func FromUser(data model.User) map[string]any {
	return map[string]any{
		"bool":       data.Bool,
		"bool_2":     data.Bool2,
		"created_at": data.CreatedAt,
		"float_32":   data.Float32,
		"float_64":   data.Float64,
		"int":        data.Int,
		"int_32":     data.Int32,
		"int_64":     data.Int64,
		"login":      FromLogin(data.Login),
		"main_group": FromGroup(data.MainGroup),
		"string":     data.String,
		"updated_at": data.UpdatedAt,
		"uuid":       data.UUID,
	}
}
func ToUser(data map[string]any) model.User {
	return model.User{
		Bool:      data["bool"].(bool),
		Bool2:     data["bool_2"].(bool),
		CreatedAt: parseTime(data["created_at"]),
		Float32:   float32(data["float_32"].(float64)),
		Float64:   data["float_64"].(float64),
		ID:        prepareID("user", data["id"]),
		Int:       int(data["int"].(float64)),
		Int32:     int32(data["int_32"].(float64)),
		Int64:     int64(data["int_64"].(float64)),
		Login:     ToLogin(data["login"].(map[string]any)),
		MainGroup: ToGroup(data["main_group"].(map[string]any)),
		String:    data["string"].(string),
		UUID:      parseUUID(data["uuid"]),
		UpdatedAt: parseTime(data["updated_at"]),
	}
}
