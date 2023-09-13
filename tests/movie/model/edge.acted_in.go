package model

import (
	"github.com/go-surreal/som"
)

type ActedIn struct {
	som.Edge

	Person *Person `som:"in"` // TODO: in and out are missing in conv!!
	Movie  *Movie  `som:"out"`
}
