package model

import (
	"github.com/go-surreal/som"
)

type Directed struct {
	som.Edge

	Person *Person `som:"in"`
	Movie  *Movie  `som:"out"`
}
