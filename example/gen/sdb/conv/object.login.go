package conv

import model "github.com/marcbinz/sdb/example/model"

type Login struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func FromLogin(data *model.Login) *Login {
	if data == nil {
		return &Login{}
	}
	return &Login{
		Password: data.Password,
		Username: data.Username,
	}
}
func ToLogin(data *Login) *model.Login {
	return &model.Login{
		Password: data.Password,
		Username: data.Username,
	}
}
