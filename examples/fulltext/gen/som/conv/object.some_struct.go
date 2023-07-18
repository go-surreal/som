// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package conv

import (
	model "github.com/marcbinz/som/examples/basic/model"
	"time"
)

type someStruct struct {
	StringPtr *string    `json:"string_ptr"`
	IntPtr    *int       `json:"int_ptr"`
	TimePtr   *time.Time `json:"time_ptr"`
	UuidPtr   *string    `json:"uuid_ptr"`
}

func fromSomeStruct(data model.SomeStruct) someStruct {
	return someStruct{
		IntPtr:    data.IntPtr,
		StringPtr: data.StringPtr,
		TimePtr:   data.TimePtr,
		UuidPtr:   uuidPtr(data.UuidPtr),
	}
}

func toSomeStruct(data someStruct) model.SomeStruct {
	return model.SomeStruct{
		IntPtr:    data.IntPtr,
		StringPtr: data.StringPtr,
		TimePtr:   data.TimePtr,
		UuidPtr:   ptrFunc(parseUUID)(data.UuidPtr),
	}
}
