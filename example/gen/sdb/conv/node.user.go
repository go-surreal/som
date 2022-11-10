package conv

import (
	model "github.com/marcbinz/sdb/example/model"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	String    string    `json:"string"`
	Int       int       `json:"int"`
	Int32     int32     `json:"int_32"`
	Int64     int64     `json:"int_64"`
	Float32   float32   `json:"float_32"`
	Float64   float64   `json:"float_64"`
	Bool      bool      `json:"bool"`
	Bool2     bool      `json:"bool_2"`
	UUID      string    `json:"uuid"`
	Login     Login     `json:"login"`
	Role      string    `json:"role"`
	Groups    []Group   `json:"groups"`
	MainGroup any       `json:"main_group"`
	Other     []string  `json:"other"`
	More      []float32 `json:"more"`
	Roles     []string  `json:"roles"`
}

func FromUser(data *model.User) *User {
	if data == nil {
		return &User{}
	}
	return &User{
		Bool:      data.Bool,
		Bool2:     data.Bool2,
		CreatedAt: data.CreatedAt,
		Float32:   data.Float32,
		Float64:   data.Float64,
		Int:       data.Int,
		Int32:     data.Int32,
		Int64:     data.Int64,
		Login:     *FromLogin(&data.Login),
		MainGroup: toGroupRecord(data.MainGroup),
		String:    data.String,
		UUID:      data.UUID.String(),
		UpdatedAt: data.UpdatedAt,
	}
}
func ToUser(data *User) *model.User {
	return &model.User{
		Bool:      data.Bool,
		Bool2:     data.Bool2,
		CreatedAt: data.CreatedAt,
		Float32:   data.Float32,
		Float64:   data.Float64,
		ID:        prepareID("user", data.ID),
		Int:       data.Int,
		Int32:     data.Int32,
		Int64:     data.Int64,
		Login:     *ToLogin(&data.Login),
		MainGroup: *fromGroupRecord(data.MainGroup),
		String:    data.String,
		UUID:      parseUUID(data.UUID),
		UpdatedAt: data.UpdatedAt,
	}
}
func fromUserRecord(data any) *model.User {
	if node, ok := data.(*User); ok {
		return ToUser(node)
	}
	return &model.User{}
}
func toUserRecord(node model.User) string {
	if node.ID == "" {
		return ""
	}
	return "user:" + node.ID
}
