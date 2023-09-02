package model

import (
	"github.com/marcbinz/som"
	"net/url"
)

type URLExample struct {
	som.Node

	SomeURL      *url.URL
	SomeOtherURL url.URL
}
