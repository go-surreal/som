// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package query

import (
	conv "github.com/go-surreal/som/tests/basic/gen/som/conv"
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

func NewFieldsLikeDBResponse(db Database, unmarshal func(buf []byte, val any) error) Builder[model.FieldsLikeDBResponse, conv.FieldsLikeDBResponse] {
	return Builder[model.FieldsLikeDBResponse, conv.FieldsLikeDBResponse]{builder[model.FieldsLikeDBResponse, conv.FieldsLikeDBResponse]{
		convFrom:  conv.FromFieldsLikeDBResponse,
		convTo:    conv.ToFieldsLikeDBResponse,
		db:        db,
		query:     lib.NewQuery[model.FieldsLikeDBResponse]("fields_like_db_response"),
		unmarshal: unmarshal,
	}}
}
