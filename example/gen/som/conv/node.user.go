package conv

import (
	"encoding/json"
	model "github.com/marcbinz/som/example/model"
	"strings"
	"time"
)

type User struct {
	ID        string       `json:"id,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
	String    string       `json:"string,omitempty"`
	Int       int          `json:"int,omitempty"`
	Int32     int32        `json:"int_32,omitempty"`
	Int64     int64        `json:"int_64,omitempty"`
	Float32   float32      `json:"float_32,omitempty"`
	Float64   float64      `json:"float_64,omitempty"`
	Bool      bool         `json:"bool,omitempty"`
	Bool2     bool         `json:"bool_2,omitempty"`
	UUID      string       `json:"uuid,omitempty"`
	Login     Login        `json:"login,omitempty"`
	Role      string       `json:"role,omitempty"`
	Groups    []GroupField `json:"groups,omitempty"`
	MainGroup GroupField   `json:"main_group,omitempty"`
	Other     []string     `json:"other,omitempty"`
	More      []float32    `json:"more,omitempty"`
	Roles     []string     `json:"roles,omitempty"`
	MyGroups  []MemberOf   `json:"my_groups,omitempty"`
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
		Groups:    mapRecords(data.Groups, toGroupField),
		ID:        buildDatabaseID("user", data.ID),
		Int:       data.Int,
		Int32:     data.Int32,
		Int64:     data.Int64,
		Login:     *FromLogin(&data.Login),
		MainGroup: toGroupField(&data.MainGroup),
		More:      data.More,
		Other:     data.Other,
		Role:      string(data.Role),
		Roles:     convertEnum[model.Role, string](data.Roles),
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
		ID:        parseDatabaseID("user", data.ID),
		Int:       data.Int,
		Int32:     data.Int32,
		Int64:     data.Int64,
		Login:     *ToLogin(&data.Login),
		MainGroup: *fromGroupField(data.MainGroup),
		More:      data.More,
		Other:     data.Other,
		Role:      model.Role(data.Role),
		Roles:     convertEnum[string, model.Role](data.Roles),
		String:    data.String,
		UUID:      parseUUID(data.UUID),
		UpdatedAt: data.UpdatedAt,
	}
}

type UserField User

func (f *UserField) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}
func (f *UserField) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\""); strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("user", raw)
		return nil
	}
	type fieldAlias UserField
	var field fieldAlias
	err := json.Unmarshal(data, &field)
	if err == nil {
		*f = UserField(field)
	}
	return err
}
func fromUserField(field UserField) *model.User {
	node := User(field)
	return ToUser(&node)
}
func toUserField(node *model.User) UserField {
	return UserField(*FromUser(node))
}