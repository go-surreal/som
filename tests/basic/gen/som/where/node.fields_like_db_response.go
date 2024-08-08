// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

var FieldsLikeDBResponse = newFieldsLikeDBResponse[model.FieldsLikeDBResponse](lib.NewKey[model.FieldsLikeDBResponse]())

func newFieldsLikeDBResponse[M any](key lib.Key[M]) fieldsLikeDbresponse[M] {
	return fieldsLikeDbresponse[M]{
		Detail: lib.NewString[M](lib.Field(key, "detail")),
		ID:     lib.NewID[M](lib.Field(key, "id"), "fields_like_db_response"),
		Result: lib.NewSliceMaker[M, string, *lib.String[M]](lib.NewString[M])(lib.Field(key, "result")),
		Status: lib.NewString[M](lib.Field(key, "status")),
		Time:   lib.NewString[M](lib.Field(key, "time")),
		key:    key,
	}
}

type fieldsLikeDbresponse[M any] struct {
	key    lib.Key[M]
	ID     *lib.ID[M]
	Time   *lib.String[M]
	Status *lib.String[M]
	Detail *lib.String[M]
	Result *lib.Slice[M, string, *lib.String[M]]
}

type fieldsLikeDbresponseEdges[M any] struct {
	lib.Filter[M]
	key lib.Key[M]
}
