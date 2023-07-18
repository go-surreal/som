package model

import (
	"github.com/marcbinz/som"
)

type Person struct {
	som.Node

	Name string
}
