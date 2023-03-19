package model

import (
	"github.com/marcbinz/som"
)

type ActedIn struct {
	som.Edge

	Person *Person `som:"in"` // TODO: in and out are missing in conv!!
	Movie  *Movie  `som:"out"`
}
