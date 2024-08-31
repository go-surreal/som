package model

import (
	"github.com/go-surreal/som"
	"net/url"
)

type URLExample struct {
	som.Node

	SomeURL      *url.URL
	SomeOtherURL url.URL
}
