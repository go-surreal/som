package model

import (
	"github.com/go-surreal/som"
)

type FieldsLikeDBResponse struct {
	som.Node

	Time   string
	Status string
	Detail string
	Result []string
}
