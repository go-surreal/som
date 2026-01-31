package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som"
)

type UUIDModel struct {
	som.CustomNode[som.UUID]

	Label string
}

type RandModel struct {
	som.CustomNode[som.Rand]

	Value int
}
