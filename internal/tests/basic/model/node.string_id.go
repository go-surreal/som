package model

import (
	"som.test/gen/som"
)

type Slug struct {
	som.Node[som.String]
	Title string
}
