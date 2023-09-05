package model

import (
	"github.com/go-surreal/som"
)

type Movie struct {
	som.Node

	Title string
}
