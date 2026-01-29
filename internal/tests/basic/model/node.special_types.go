package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som"
)

type SpecialTypes struct {
	som.Node
	som.OptimisticLock
	som.SoftDelete

	Name string
}

type SpecialRelation struct {
	som.Node
	som.SoftDelete

	Title   string
	Author  *SpecialTypes
	Authors []*SpecialTypes
}
