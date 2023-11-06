// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	"encoding/json"
	som "github.com/go-surreal/som"
	model "github.com/go-surreal/som/tests/basic/model"
	uuid "github.com/google/uuid"
	"strings"
	"time"
)

type AllFieldTypes struct {
	ID                string         `json:"id,omitempty"`
	CreatedAt         *time.Time     `json:"created_at,omitempty"`
	UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
	String            string         `json:"string"`
	StringPtr         *string        `json:"string_ptr"`
	Other             []string       `json:"other"`
	StringPtrSlice    []*string      `json:"string_ptr_slice"`
	StringSlicePtr    *[]string      `json:"string_slice_ptr"`
	Int               int            `json:"int"`
	IntPtr            *int           `json:"int_ptr"`
	Int8              int8           `json:"int_8"`
	Int8Ptr           *int8          `json:"int_8_ptr"`
	Int16             int16          `json:"int_16"`
	Int16Ptr          *int16         `json:"int_16_ptr"`
	Int32             int32          `json:"int_32"`
	Int32Ptr          *int32         `json:"int_32_ptr"`
	Int64             int64          `json:"int_64"`
	Int64Ptr          *int64         `json:"int_64_ptr"`
	Uint8             uint8          `json:"uint_8"`
	Uint8Ptr          *uint8         `json:"uint_8_ptr"`
	Uint16            uint16         `json:"uint_16"`
	Uint16Ptr         *uint16        `json:"uint_16_ptr"`
	Uint32            uint32         `json:"uint_32"`
	Uint32Ptr         *uint32        `json:"uint_32_ptr"`
	Float32           float32        `json:"float_32"`
	More              []float32      `json:"more"`
	Float64           float64        `json:"float_64"`
	Rune              rune           `json:"rune"`
	Bool              bool           `json:"bool"`
	BoolPtr           *bool          `json:"bool_ptr"`
	Time              time.Time      `json:"time"`
	TimePtr           *time.Time     `json:"time_ptr"`
	UUID              uuid.UUID      `json:"uuid"`
	UUIDPtr           *uuid.UUID     `json:"uuid_ptr"`
	Role              model.Role     `json:"role"`
	EnumPtr           *model.Role    `json:"enum_ptr"`
	Roles             []model.Role   `json:"roles"`
	EnumPtrSlice      []*model.Role  `json:"enum_ptr_slice"`
	EnumPtrSlicePtr   *[]*model.Role `json:"enum_ptr_slice_ptr"`
	Login             login          `json:"login"`
	StructPtr         *someStruct    `json:"struct_ptr"`
	StructSlice       []someStruct   `json:"struct_slice"`
	StructPtrSlice    []*someStruct  `json:"struct_ptr_slice"`
	StructPtrSlicePtr *[]*someStruct `json:"struct_ptr_slice_ptr"`
	MainGroup         *groupLink     `json:"main_group"`
	MainGroupPtr      *groupLink     `json:"main_group_ptr"`
	Groups            []*groupLink   `json:"groups"`
	NodePtrSlice      []*groupLink   `json:"node_ptr_slice"`
	NodePtrSlicePtr   *[]*groupLink  `json:"node_ptr_slice_ptr"`
	MemberOf          []GroupMember  `json:"member_of,omitempty"`
	SliceSlice        [][]string     `json:"slice_slice"`
}

func FromAllFieldTypes(data *model.AllFieldTypes) *AllFieldTypes {
	if data == nil {
		return nil
	}
	return &AllFieldTypes{
		Bool:              data.Bool,
		BoolPtr:           data.BoolPtr,
		EnumPtr:           data.EnumPtr,
		EnumPtrSlice:      data.EnumPtrSlice,
		EnumPtrSlicePtr:   data.EnumPtrSlicePtr,
		Float32:           data.Float32,
		Float64:           data.Float64,
		Groups:            mapSlice(data.Groups, toGroupLink),
		Int:               data.Int,
		Int16:             data.Int16,
		Int16Ptr:          data.Int16Ptr,
		Int32:             data.Int32,
		Int32Ptr:          data.Int32Ptr,
		Int64:             data.Int64,
		Int64Ptr:          data.Int64Ptr,
		Int8:              data.Int8,
		Int8Ptr:           data.Int8Ptr,
		IntPtr:            data.IntPtr,
		Login:             noPtrFunc(fromLogin)(data.Login),
		MainGroup:         toGroupLink(data.MainGroup),
		MainGroupPtr:      toGroupLinkPtr(data.MainGroupPtr),
		More:              data.More,
		NodePtrSlice:      mapSlice(data.NodePtrSlice, toGroupLinkPtr),
		NodePtrSlicePtr:   mapSlicePtr(data.NodePtrSlicePtr, toGroupLinkPtr),
		Other:             data.Other,
		Role:              data.Role,
		Roles:             data.Roles,
		Rune:              data.Rune,
		SliceSlice:        data.SliceSlice,
		String:            data.String,
		StringPtr:         data.StringPtr,
		StringPtrSlice:    data.StringPtrSlice,
		StringSlicePtr:    data.StringSlicePtr,
		StructPtr:         fromSomeStruct(data.StructPtr),
		StructPtrSlice:    mapSlice(data.StructPtrSlice, fromSomeStruct),
		StructPtrSlicePtr: mapSlicePtr(data.StructPtrSlicePtr, fromSomeStruct),
		StructSlice:       mapSlice(data.StructSlice, noPtrFunc(fromSomeStruct)),
		Time:              data.Time,
		TimePtr:           data.TimePtr,
		UUID:              data.UUID,
		UUIDPtr:           data.UUIDPtr,
		Uint16:            data.Uint16,
		Uint16Ptr:         data.Uint16Ptr,
		Uint32:            data.Uint32,
		Uint32Ptr:         data.Uint32Ptr,
		Uint8:             data.Uint8,
		Uint8Ptr:          data.Uint8Ptr,
	}
}

func ToAllFieldTypes(data *AllFieldTypes) *model.AllFieldTypes {
	if data == nil {
		return nil
	}
	return &model.AllFieldTypes{
		Bool:              data.Bool,
		BoolPtr:           data.BoolPtr,
		EnumPtr:           data.EnumPtr,
		EnumPtrSlice:      data.EnumPtrSlice,
		EnumPtrSlicePtr:   data.EnumPtrSlicePtr,
		Float32:           data.Float32,
		Float64:           data.Float64,
		Groups:            mapSlice(data.Groups, fromGroupLink),
		Int:               data.Int,
		Int16:             data.Int16,
		Int16Ptr:          data.Int16Ptr,
		Int32:             data.Int32,
		Int32Ptr:          data.Int32Ptr,
		Int64:             data.Int64,
		Int64Ptr:          data.Int64Ptr,
		Int8:              data.Int8,
		Int8Ptr:           data.Int8Ptr,
		IntPtr:            data.IntPtr,
		Login:             noPtrFunc(toLogin)(data.Login),
		MainGroup:         fromGroupLink(data.MainGroup),
		MainGroupPtr:      fromGroupLinkPtr(data.MainGroupPtr),
		MemberOf:          mapSlice(data.MemberOf, noPtrFunc(ToGroupMember)),
		More:              data.More,
		Node:              som.NewNode(parseDatabaseID("all_field_types", data.ID)),
		NodePtrSlice:      mapSlice(data.NodePtrSlice, fromGroupLinkPtr),
		NodePtrSlicePtr:   mapSlicePtr(data.NodePtrSlicePtr, fromGroupLinkPtr),
		Other:             data.Other,
		Role:              data.Role,
		Roles:             data.Roles,
		Rune:              data.Rune,
		SliceSlice:        data.SliceSlice,
		String:            data.String,
		StringPtr:         data.StringPtr,
		StringPtrSlice:    data.StringPtrSlice,
		StringSlicePtr:    data.StringSlicePtr,
		StructPtr:         toSomeStruct(data.StructPtr),
		StructPtrSlice:    mapSlice(data.StructPtrSlice, toSomeStruct),
		StructPtrSlicePtr: mapSlicePtr(data.StructPtrSlicePtr, toSomeStruct),
		StructSlice:       mapSlice(data.StructSlice, noPtrFunc(toSomeStruct)),
		Time:              data.Time,
		TimePtr:           data.TimePtr,
		Timestamps:        som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
		UUID:              data.UUID,
		UUIDPtr:           data.UUIDPtr,
		Uint16:            data.Uint16,
		Uint16Ptr:         data.Uint16Ptr,
		Uint32:            data.Uint32,
		Uint32Ptr:         data.Uint32Ptr,
		Uint8:             data.Uint8,
		Uint8Ptr:          data.Uint8Ptr,
	}
}

type allFieldTypesLink struct {
	AllFieldTypes
	ID string
}

func (f *allFieldTypesLink) MarshalJSON() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f.ID)
}

func (f *allFieldTypesLink) UnmarshalJSON(data []byte) error {
	raw := string(data)
	if strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"") {
		raw = raw[1 : len(raw)-1]
		f.ID = parseDatabaseID("all_field_types", raw)
		return nil
	}
	type alias allFieldTypesLink
	var link alias
	err := json.Unmarshal(data, &link)
	if err == nil {
		*f = allFieldTypesLink(link)
	}
	return err
}

func fromAllFieldTypesLink(link *allFieldTypesLink) model.AllFieldTypes {
	if link == nil {
		return model.AllFieldTypes{}
	}
	res := AllFieldTypes(link.AllFieldTypes)
	out := ToAllFieldTypes(&res)
	return *out
}

func fromAllFieldTypesLinkPtr(link *allFieldTypesLink) *model.AllFieldTypes {
	if link == nil {
		return nil
	}
	res := AllFieldTypes(link.AllFieldTypes)
	return ToAllFieldTypes(&res)
}

func toAllFieldTypesLink(node model.AllFieldTypes) *allFieldTypesLink {
	if node.ID() == "" {
		return nil
	}
	link := allFieldTypesLink{AllFieldTypes: *FromAllFieldTypes(&node), ID: buildDatabaseID("all_field_types", node.ID())}
	return &link
}

func toAllFieldTypesLinkPtr(node *model.AllFieldTypes) *allFieldTypesLink {
	if node == nil || node.ID() == "" {
		return nil
	}
	link := allFieldTypesLink{AllFieldTypes: *FromAllFieldTypes(node), ID: buildDatabaseID("all_field_types", node.ID())}
	return &link
}
