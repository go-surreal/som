package conv

import model "github.com/marcbinz/sdb/example/model"

func FromLogin(data model.Login) map[string]any {
	return map[string]any{
		"password": data.Password,
		"username": data.Username,
	}
}
func ToLogin(data map[string]any) model.Login {
	return model.Login{
		Password: data["password"].(string),
		Username: data["username"].(string),
	}
}
