package model

import "github.com/go-surreal/som/tests/basic/gen/som"

type SoftDeleteUser struct {
	som.Node
	som.SoftDelete
	Name string
}

type SoftDeleteComplete struct {
	som.Node
	som.Timestamps
	som.OptimisticLock
	som.SoftDelete
	Name string
}
