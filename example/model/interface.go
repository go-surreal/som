package model

import (
	"github.com/marcbinz/som"
)

type Shape interface {
	som.Node
	typeShape()
}

// https://stackoverflow.com/questions/28033277/decoding-generic-json-objects-to-one-of-many-formats
var _ Shape = &Square{}
var _ Shape = &Triangle{}
var _ Shape = &Circle{}

type Square struct {
	som.Node
}

func (*Square) typeShape() {}

type Triangle struct {
	som.Node
}

func (*Triangle) typeShape() {}

type Circle struct {
	som.Node
}

func (*Circle) typeShape() {}
