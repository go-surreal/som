// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package conv

import (
	v2 "github.com/fxamacker/cbor/v2"
	sdbc "github.com/go-surreal/sdbc"
	som "github.com/go-surreal/som"
	types "github.com/go-surreal/som/tests/basic/gen/som/internal/types"
	model "github.com/go-surreal/som/tests/basic/model"
)

type AllFieldTypes struct {
	ID                 *sdbc.ID          `cbor:"id,omitempty"`
	CreatedAt          *sdbc.DateTime    `cbor:"created_at,omitempty"`
	UpdatedAt          *sdbc.DateTime    `cbor:"updated_at,omitempty"`
	String             string            `cbor:"string"`
	StringPtr          *string           `cbor:"string_ptr"`
	Other              []string          `cbor:"other"`
	StringPtrSlice     []*string         `cbor:"string_ptr_slice"`
	StringSlicePtr     *[]string         `cbor:"string_slice_ptr"`
	Int                int               `cbor:"int"`
	IntPtr             *int              `cbor:"int_ptr"`
	IntSlice           []int             `cbor:"int_slice"`
	IntPtrSlice        []*int            `cbor:"int_ptr_slice"`
	IntSlicePtr        *[]int            `cbor:"int_slice_ptr"`
	IntPtrSlicePtr     *[]*int           `cbor:"int_ptr_slice_ptr"`
	Int8               int8              `cbor:"int_8"`
	Int8Ptr            *int8             `cbor:"int_8_ptr"`
	Int16              int16             `cbor:"int_16"`
	Int16Ptr           *int16            `cbor:"int_16_ptr"`
	Int32              int32             `cbor:"int_32"`
	Int32Ptr           *int32            `cbor:"int_32_ptr"`
	Int64              int64             `cbor:"int_64"`
	Int64Ptr           *int64            `cbor:"int_64_ptr"`
	Uint8              uint8             `cbor:"uint_8"`
	Uint8Ptr           *uint8            `cbor:"uint_8_ptr"`
	Uint16             uint16            `cbor:"uint_16"`
	Uint16Ptr          *uint16           `cbor:"uint_16_ptr"`
	Uint32             uint32            `cbor:"uint_32"`
	Uint32Ptr          *uint32           `cbor:"uint_32_ptr"`
	Float32            float32           `cbor:"float_32"`
	Float32Slice       []float32         `cbor:"float_32_slice"`
	Float32SlicePtr    *[]float32        `cbor:"float_32_slice_ptr"`
	Float32PtrSlice    []*float32        `cbor:"float_32_ptr_slice"`
	Float32PtrSlicePtr *[]*float32       `cbor:"float_32_ptr_slice_ptr"`
	Float64            float64           `cbor:"float_64"`
	Rune               rune              `cbor:"rune"`
	RuneSlice          []rune            `cbor:"rune_slice"`
	Bool               bool              `cbor:"bool"`
	BoolPtr            *bool             `cbor:"bool_ptr"`
	BoolSlice          []bool            `cbor:"bool_slice"`
	Time               sdbc.DateTime     `cbor:"time"`
	TimePtr            *sdbc.DateTime    `cbor:"time_ptr"`
	TimeNil            *sdbc.DateTime    `cbor:"time_nil"`
	TimeSlice          []sdbc.DateTime   `cbor:"time_slice"`
	TimeSliceSlice     [][]sdbc.DateTime `cbor:"time_slice_slice"`
	Duration           sdbc.Duration     `cbor:"duration"`
	DurationPtr        *sdbc.Duration    `cbor:"duration_ptr"`
	DurationNil        *sdbc.Duration    `cbor:"duration_nil"`
	DurationSlice      []sdbc.Duration   `cbor:"duration_slice"`
	UUID               types.UUID        `cbor:"uuid"`
	UUIDPtr            *types.UUID       `cbor:"uuid_ptr"`
	UUIDNil            *types.UUID       `cbor:"uuid_nil"`
	UUIDSlice          []types.UUID      `cbor:"uuid_slice"`
	URL                string            `cbor:"url"`
	URLPtr             *string           `cbor:"url_ptr"`
	URLNil             *string           `cbor:"url_nil"`
	URLSlice           []string          `cbor:"url_slice"`
	Role               model.Role        `cbor:"role"`
	EnumPtr            *model.Role       `cbor:"enum_ptr"`
	Roles              []model.Role      `cbor:"roles"`
	EnumPtrSlice       []*model.Role     `cbor:"enum_ptr_slice"`
	EnumPtrSlicePtr    *[]*model.Role    `cbor:"enum_ptr_slice_ptr"`
	Login              login             `cbor:"login"`
	StructPtr          *someStruct       `cbor:"struct_ptr"`
	StructSlice        []someStruct      `cbor:"struct_slice"`
	StructPtrSlice     []*someStruct     `cbor:"struct_ptr_slice"`
	StructPtrSlicePtr  *[]*someStruct    `cbor:"struct_ptr_slice_ptr"`
	MainGroup          *groupLink        `cbor:"main_group"`
	MainGroupPtr       *groupLink        `cbor:"main_group_ptr"`
	Groups             []*groupLink      `cbor:"groups"`
	GroupsSlice        [][]*groupLink    `cbor:"groups_slice"`
	NodePtrSlice       []*groupLink      `cbor:"node_ptr_slice"`
	NodePtrSlicePtr    *[]*groupLink     `cbor:"node_ptr_slice_ptr"`
	MemberOf           []GroupMember     `cbor:"member_of,omitempty"`
	SliceSlice         [][]string        `cbor:"slice_slice"`
	SliceSliceSlice    [][][]string      `cbor:"slice_slice_slice"`
	SliceSliceSlice2   [][][]someStruct  `cbor:"slice_slice_slice_2"`
	Byte               byte              `cbor:"byte"`
	BytePtr            *byte             `cbor:"byte_ptr"`
	ByteSlice          []byte            `cbor:"byte_slice"`
	ByteSlicePtr       *[]byte           `cbor:"byte_slice_ptr"`
}

func FromAllFieldTypes(data *model.AllFieldTypes) *AllFieldTypes {
	if data == nil {
		return nil
	}
	return &AllFieldTypes{
		Bool:               data.Bool,
		BoolPtr:            data.BoolPtr,
		BoolSlice:          data.BoolSlice,
		Byte:               data.Byte,
		BytePtr:            data.BytePtr,
		ByteSlice:          data.ByteSlice,
		ByteSlicePtr:       data.ByteSlicePtr,
		Duration:           fromDuration(data.Duration),
		DurationNil:        fromDurationPtr(data.DurationNil),
		DurationPtr:        fromDurationPtr(data.DurationPtr),
		DurationSlice:      mapSliceFn(fromDuration)(data.DurationSlice),
		EnumPtr:            data.EnumPtr,
		Float32:            data.Float32,
		Float32PtrSlice:    data.Float32PtrSlice,
		Float32PtrSlicePtr: data.Float32PtrSlicePtr,
		Float32Slice:       data.Float32Slice,
		Float32SlicePtr:    data.Float32SlicePtr,
		Float64:            data.Float64,
		Groups:             mapSliceFn(toGroupLink)(data.Groups),
		GroupsSlice:        mapSliceFn(mapSliceFn(toGroupLink))(data.GroupsSlice),
		Int:                data.Int,
		Int16:              data.Int16,
		Int16Ptr:           data.Int16Ptr,
		Int32:              data.Int32,
		Int32Ptr:           data.Int32Ptr,
		Int64:              data.Int64,
		Int64Ptr:           data.Int64Ptr,
		Int8:               data.Int8,
		Int8Ptr:            data.Int8Ptr,
		IntPtr:             data.IntPtr,
		IntPtrSlice:        data.IntPtrSlice,
		IntPtrSlicePtr:     data.IntPtrSlicePtr,
		IntSlice:           data.IntSlice,
		IntSlicePtr:        data.IntSlicePtr,
		Login:              noPtrFunc(fromLogin)(data.Login),
		MainGroup:          toGroupLink(data.MainGroup),
		MainGroupPtr:       toGroupLinkPtr(data.MainGroupPtr),
		NodePtrSlice:       mapSliceFn(toGroupLinkPtr)(data.NodePtrSlice),
		NodePtrSlicePtr:    mapSliceFnPtr(toGroupLinkPtr)(data.NodePtrSlicePtr),
		Other:              data.Other,
		Role:               data.Role,
		Rune:               data.Rune,
		RuneSlice:          data.RuneSlice,
		SliceSlice:         data.SliceSlice,
		SliceSliceSlice:    data.SliceSliceSlice,
		SliceSliceSlice2:   mapSliceFn(mapSliceFn(mapSliceFn(noPtrFunc(fromSomeStruct))))(data.SliceSliceSlice2),
		String:             data.String,
		StringPtr:          data.StringPtr,
		StringPtrSlice:     data.StringPtrSlice,
		StringSlicePtr:     data.StringSlicePtr,
		StructPtr:          fromSomeStruct(data.StructPtr),
		StructPtrSlice:     mapSliceFn(fromSomeStruct)(data.StructPtrSlice),
		StructPtrSlicePtr:  mapSliceFnPtr(fromSomeStruct)(data.StructPtrSlicePtr),
		StructSlice:        mapSliceFn(noPtrFunc(fromSomeStruct))(data.StructSlice),
		Time:               fromTime(data.Time),
		TimeNil:            fromTimePtr(data.TimeNil),
		TimePtr:            fromTimePtr(data.TimePtr),
		TimeSlice:          mapSliceFn(fromTime)(data.TimeSlice),
		TimeSliceSlice:     mapSliceFn(mapSliceFn(fromTime))(data.TimeSliceSlice),
		URL:                fromURL(data.URL),
		URLNil:             fromURLPtr(data.URLNil),
		URLPtr:             fromURLPtr(data.URLPtr),
		URLSlice:           mapSliceFn(fromURL)(data.URLSlice),
		UUID:               fromUUID(data.UUID),
		UUIDNil:            fromUUIDPtr(data.UUIDNil),
		UUIDPtr:            fromUUIDPtr(data.UUIDPtr),
		UUIDSlice:          mapSliceFn(fromUUID)(data.UUIDSlice),
		Uint16:             data.Uint16,
		Uint16Ptr:          data.Uint16Ptr,
		Uint32:             data.Uint32,
		Uint32Ptr:          data.Uint32Ptr,
		Uint8:              data.Uint8,
		Uint8Ptr:           data.Uint8Ptr,
	}
}

func ToAllFieldTypes(data *AllFieldTypes) *model.AllFieldTypes {
	if data == nil {
		return nil
	}
	return &model.AllFieldTypes{
		Bool:               data.Bool,
		BoolPtr:            data.BoolPtr,
		BoolSlice:          data.BoolSlice,
		Byte:               data.Byte,
		BytePtr:            data.BytePtr,
		ByteSlice:          data.ByteSlice,
		ByteSlicePtr:       data.ByteSlicePtr,
		Duration:           toDuration(data.Duration),
		DurationNil:        toDurationPtr(data.DurationNil),
		DurationPtr:        toDurationPtr(data.DurationPtr),
		DurationSlice:      mapSliceFn(toDuration)(data.DurationSlice),
		EnumPtr:            data.EnumPtr,
		Float32:            data.Float32,
		Float32PtrSlice:    data.Float32PtrSlice,
		Float32PtrSlicePtr: data.Float32PtrSlicePtr,
		Float32Slice:       data.Float32Slice,
		Float32SlicePtr:    data.Float32SlicePtr,
		Float64:            data.Float64,
		Groups:             mapSliceFn(fromGroupLink)(data.Groups),
		GroupsSlice:        mapSliceFn(mapSliceFn(fromGroupLink))(data.GroupsSlice),
		Int:                data.Int,
		Int16:              data.Int16,
		Int16Ptr:           data.Int16Ptr,
		Int32:              data.Int32,
		Int32Ptr:           data.Int32Ptr,
		Int64:              data.Int64,
		Int64Ptr:           data.Int64Ptr,
		Int8:               data.Int8,
		Int8Ptr:            data.Int8Ptr,
		IntPtr:             data.IntPtr,
		IntPtrSlice:        data.IntPtrSlice,
		IntPtrSlicePtr:     data.IntPtrSlicePtr,
		IntSlice:           data.IntSlice,
		IntSlicePtr:        data.IntSlicePtr,
		Login:              noPtrFunc(toLogin)(data.Login),
		MainGroup:          fromGroupLink(data.MainGroup),
		MainGroupPtr:       fromGroupLinkPtr(data.MainGroupPtr),
		MemberOf:           mapSliceFn(noPtrFunc(ToGroupMember))(data.MemberOf),
		Node:               som.NewNode(data.ID),
		NodePtrSlice:       mapSliceFn(fromGroupLinkPtr)(data.NodePtrSlice),
		NodePtrSlicePtr:    mapSliceFnPtr(fromGroupLinkPtr)(data.NodePtrSlicePtr),
		Other:              data.Other,
		Role:               data.Role,
		Rune:               data.Rune,
		RuneSlice:          data.RuneSlice,
		SliceSlice:         data.SliceSlice,
		SliceSliceSlice:    data.SliceSliceSlice,
		SliceSliceSlice2:   mapSliceFn(mapSliceFn(mapSliceFn(noPtrFunc(toSomeStruct))))(data.SliceSliceSlice2),
		String:             data.String,
		StringPtr:          data.StringPtr,
		StringPtrSlice:     data.StringPtrSlice,
		StringSlicePtr:     data.StringSlicePtr,
		StructPtr:          toSomeStruct(data.StructPtr),
		StructPtrSlice:     mapSliceFn(toSomeStruct)(data.StructPtrSlice),
		StructPtrSlicePtr:  mapSliceFnPtr(toSomeStruct)(data.StructPtrSlicePtr),
		StructSlice:        mapSliceFn(noPtrFunc(toSomeStruct))(data.StructSlice),
		Time:               toTime(data.Time),
		TimeNil:            toTimePtr(data.TimeNil),
		TimePtr:            toTimePtr(data.TimePtr),
		TimeSlice:          mapSliceFn(toTime)(data.TimeSlice),
		TimeSliceSlice:     mapSliceFn(mapSliceFn(toTime))(data.TimeSliceSlice),
		Timestamps:         som.NewTimestamps(data.CreatedAt, data.UpdatedAt),
		URL:                toURL(data.URL),
		URLNil:             toURLPtr(data.URLNil),
		URLPtr:             toURLPtr(data.URLPtr),
		URLSlice:           mapSliceFn(toURL)(data.URLSlice),
		UUID:               toUUID(data.UUID),
		UUIDNil:            toUUIDPtr(data.UUIDNil),
		UUIDPtr:            toUUIDPtr(data.UUIDPtr),
		UUIDSlice:          mapSliceFn(toUUID)(data.UUIDSlice),
		Uint16:             data.Uint16,
		Uint16Ptr:          data.Uint16Ptr,
		Uint32:             data.Uint32,
		Uint32Ptr:          data.Uint32Ptr,
		Uint8:              data.Uint8,
		Uint8Ptr:           data.Uint8Ptr,
	}
}

type allFieldTypesLink struct {
	AllFieldTypes
	ID *sdbc.ID
}

func (f *allFieldTypesLink) MarshalCBOR() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	return v2.Marshal(f.ID)
}

func (f *allFieldTypesLink) UnmarshalCBOR(data []byte) error {
	if err := v2.Unmarshal(data, &f.ID); err == nil {
		return nil
	}
	type alias allFieldTypesLink
	var link alias
	err := v2.Unmarshal(data, &link)
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
	if node.ID() == nil {
		return nil
	}
	link := allFieldTypesLink{AllFieldTypes: *FromAllFieldTypes(&node), ID: node.ID()}
	return &link
}

func toAllFieldTypesLinkPtr(node *model.AllFieldTypes) *allFieldTypesLink {
	if node == nil || node.ID() == nil {
		return nil
	}
	link := allFieldTypesLink{AllFieldTypes: *FromAllFieldTypes(node), ID: node.ID()}
	return &link
}
