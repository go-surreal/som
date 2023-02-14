// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import (
	"encoding/json"
	som "github.com/marcbinz/som"
	model "github.com/marcbinz/som/examples/basic/model"
	"strings"
	"time"
)

type User struct {
	ID                string         `json:"id,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	String            string         `json:"string"`
	Int               int            `json:"int"`
	Int32             int32          `json:"int_32"`
	Int64             int64          `json:"int_64"`
	Float32           float32        `json:"float_32"`
	Float64           float64        `json:"float_64"`
	Bool              bool           `json:"bool"`
	Bool2             bool           `json:"bool_2"`
	UUID              string         `json:"uuid"`
	Login             login          `json:"login"`
	Role              string         `json:"role"`
	Groups            []*groupLink   `json:"groups"`
	MainGroup         *groupLink     `json:"main_group"`
	MainGroupPtr      *groupLink     `json:"main_group_ptr"`
	Other             []string       `json:"other"`
	More              []float32      `json:"more"`
	Roles             []string       `json:"roles"`
	MemberOf          []GroupMember  `json:"member_of,omitempty"`
	StringPtr         *string        `json:"string_ptr"`
	IntPtr            *int           `json:"int_ptr"`
	TimePtr           *time.Time     `json:"time_ptr"`
	UuidPtr           *string        `json:"uuid_ptr"`
	StructPtr         *someStruct    `json:"struct_ptr"`
	StringPtrSlice    []*string      `json:"string_ptr_slice"`
	StringSlicePtr    *[]string      `json:"string_slice_ptr"`
	StructPtrSlice    []*someStruct  `json:"struct_ptr_slice"`
	StructPtrSlicePtr *[]*someStruct `json:"struct_ptr_slice_ptr"`
	EnumPtrSlice      []*string      `json:"enum_ptr_slice"`
	NodePtrSlice      []*groupLink   `json:"node_ptr_slice"`
	NodePtrSlicePtr   *[]*groupLink  `json:"node_ptr_slice_ptr"`
	SliceSlice        [][]string     `json:"slice_slice"`
}

func FromUser(data model.User) User {
	return User{
		Bool:              data.Bool,
		Bool2:             data.Bool2,
		EnumPtrSlice:      mapSlice(data.EnumPtrSlice, ptrFunc(mapEnum[model.Role, string])),
		Float32:           data.Float32,
		Float64:           data.Float64,
		Groups:            mapSlice(data.Groups, toGroupLink),
		Int:               data.Int,
		Int32:             data.Int32,
		Int64:             data.Int64,
		IntPtr:            data.IntPtr,
		Login:             fromLogin(data.Login),
		MainGroup:         toGroupLink(data.MainGroup),
		MainGroupPtr:      toGroupLinkPtr(data.MainGroupPtr),
		More:              data.More,
		NodePtrSlice:      mapSlice(data.NodePtrSlice, toGroupLinkPtr),
		NodePtrSlicePtr:   mapSlicePtr(data.NodePtrSlicePtr, toGroupLinkPtr),
		Other:             data.Other,
		Role:              string(data.Role),
		Roles:             mapSlice(data.Roles, mapEnum[model.Role, string]),
		SliceSlice:        data.SliceSlice,
		String:            data.String,
		StringPtr:         data.StringPtr,
		StringPtrSlice:    data.StringPtrSlice,
		StringSlicePtr:    data.StringSlicePtr,
		StructPtr:         ptrFunc(fromSomeStruct)(data.StructPtr),
		StructPtrSlice:    mapPtrSlice(data.StructPtrSlice, fromSomeStruct),
		StructPtrSlicePtr: mapPtrSlicePtr(data.StructPtrSlicePtr, fromSomeStruct),
		TimePtr:           data.TimePtr,
		UUID:              data.UUID.String(),
		UuidPtr:           uuidPtr(data.UuidPtr),
	}
}

func ToUser(data User) model.User {
	return model.User{
		Bool:              data.Bool,
		Bool2:             data.Bool2,
		EnumPtrSlice:      mapSlice(data.EnumPtrSlice, ptrFunc(mapEnum[string, model.Role])),
		Float32:           data.Float32,
		Float64:           data.Float64,
		Groups:            mapSlice(data.Groups, fromGroupLink),
		Int:               data.Int,
		Int32:             data.Int32,
		Int64:             data.Int64,
		IntPtr:            data.IntPtr,
		Login:             toLogin(data.Login),
		MainGroup:         fromGroupLink(data.MainGroup),
		MainGroupPtr:      fromGroupLinkPtr(data.MainGroupPtr),
		MemberOf:          mapSlice(data.MemberOf, ToGroupMember),
		More:              data.More,
		Node:              som.NewNode(parseDatabaseID("user", data.ID)),
		NodePtrSlice:      mapSlice(data.NodePtrSlice, fromGroupLinkPtr),
		NodePtrSlicePtr:   mapSlicePtr(data.NodePtrSlicePtr, fromGroupLinkPtr),
		Other:             data.Other,
		Role:              model.Role(data.Role),
		Roles:             mapSlice(data.Roles, mapEnum[string, model.Role]),
		SliceSlice:        data.SliceSlice,
		String:            data.String,
		StringPtr:         data.StringPtr,
		StringPtrSlice:    data.StringPtrSlice,
		StringSlicePtr:    data.StringSlicePtr,
		StructPtr:         ptrFunc(toSomeStruct)(data.StructPtr),
		StructPtrSlice:    mapPtrSlice(data.StructPtrSlice, toSomeStruct),
		StructPtrSlicePtr: mapPtrSlicePtr(data.StructPtrSlicePtr, toSomeStruct),
		TimePtr:           data.TimePtr,
		Timestamps:        som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
		UUID:              parseUUID(data.UUID),
		UuidPtr:           ptrFunc(parseUUID)(data.UuidPtr),
	}
}

type userLink struct {
	User
	ID string
}

func (f *userLink) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}

func (f *userLink) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("user", raw)
		return nil
	}
	type alias userLink
	var link alias
	err := json.Unmarshal(data, &link)
	if err == nil {
		*f = userLink(link)
	}
	return err
}

func fromUserLink(link *userLink) model.User {
	if link == nil {
		return model.User{}
	}
	return ToUser(User(link.User))
}

func fromUserLinkPtr(link *userLink) *model.User {
	if link == nil {
		return nil
	}
	node := ToUser(User(link.User))
	return &node
}

func toUserLink(node model.User) *userLink {
	if node.ID() == "" {
		return nil
	}
	link := userLink{User: FromUser(node), ID: buildDatabaseID("user", node.ID())}
	return &link
}

func toUserLinkPtr(node *model.User) *userLink {
	if node == nil || node.ID() == "" {
		return nil
	}
	link := userLink{User: FromUser(*node), ID: buildDatabaseID("user", node.ID())}
	return &link
}
