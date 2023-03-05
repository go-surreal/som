package model

import (
	"github.com/marcbinz/som"
)

type Movie struct {
	som.Node

	Title string
}
