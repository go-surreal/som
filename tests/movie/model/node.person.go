package model

import (
	"github.com/go-surreal/som"
)

type Person struct {
	som.Node

	Name string
}
