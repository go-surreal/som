package model

import (
	"github.com/marcbinz/som"
)

type Directed struct {
	som.Edge

	Person *Person `som:"in"`
	Movie  *Movie  `som:"out"`
}
