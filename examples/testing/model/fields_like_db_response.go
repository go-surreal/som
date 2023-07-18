package model

import (
	"github.com/marcbinz/som"
)

type FieldsLikeDBResponse struct {
	som.Node

	Time   string
	Status string
	Detail string
	Result []string
}
