package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som"
)

type SpecialTypes struct {
	som.CustomNode[som.UUID]
	som.OptimisticLock
	som.SoftDelete

	Name string
}

type SpecialRelation struct {
	som.CustomNode[som.Rand]
	som.SoftDelete

	Title   string
	Author  *SpecialTypes
	Authors []*SpecialTypes
}
