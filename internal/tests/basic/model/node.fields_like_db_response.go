package model

import "github.com/go-surreal/som/tests/basic/gen/som/sombase"

type FieldsLikeDBResponse struct {
	sombase.Node

	Time   string
	Status string
	Detail string
	Result []string
}
