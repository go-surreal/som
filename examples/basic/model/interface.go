package model

import (
	"github.com/marcbinz/som"
)

type Shape interface {
	som.Union
	shape
}

type shape interface {
	typeShape()
}

// https://stackoverflow.com/questions/28033277/decoding-generic-json-objects-to-one-of-many-formats
var _ Shape = &Square{}
var _ Shape = &Triangle{}
var _ Shape = &Circle{}

type Square struct {
	som.Node
	shape
}

type Triangle struct {
	som.Node
	shape
}

type Circle struct {
	som.Node
	shape
}
