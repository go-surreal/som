package model

import (
	"github.com/go-surreal/som"
)

type VersionedExample struct {
	som.Node
	som.Versioned

	SomeValue string
}
