package model

import (
	"math"

	"som.test/gen/som"
)

// Shape is a union of the node models a Drawing can be about. A field typed as
// Shape becomes a multi-table record link in the database.
//
// Next to the som.Union embed a union may declare methods of its own, which
// every member then has to implement.
type Shape interface {
	som.Union[Shape]

	Area() float64
}

// Colored is an open union: it declares no marker method, so a node can be a
// member of it next to the sealed Shape.
type Colored interface {
	som.OpenUnion

	Color() string
}

type Square struct {
	som.Node[som.ULID]
	som.Part[Shape]

	_ som.Part[Colored]

	Label string
	Size  float64
}

func (s Square) Color() string {
	return "red"
}

func (s Square) Area() float64 {
	return s.Size * s.Size
}

type Circle struct {
	som.Node[som.UUID]
	som.Part[Shape]

	_ som.Part[Colored]

	Label  string
	Radius float64
}

func (c Circle) Color() string {
	return "blue"
}

// Area has a pointer receiver on purpose, so that only *Circle implements Shape.
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Drawing struct {
	som.Node[som.ULID]
	som.Timestamps

	Name string

	// Subject links to either a Square or a Circle record.
	Subject Shape

	// Accent links to a member of the open union Colored.
	Accent Colored

	Annotations []Annotates
}

// Annotates relates a drawing to any shape, so its outgoing end is a
// multi-table relation.
type Annotates struct {
	som.Edge

	Drawing Drawing `som:"in"`
	Shape   Shape   `som:"out"`

	Note string
}
