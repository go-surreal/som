package model

import (
	"som.test/gen/som"
)

type SpecialTypes struct {
	som.Node[som.UUID]
	som.OptimisticLock
	som.SoftDelete

	Name string
}

type SpecialRelation struct {
	som.Node[som.Rand]
	som.SoftDelete

	Title   string
	Author  *SpecialTypes
	Authors []*SpecialTypes
}
