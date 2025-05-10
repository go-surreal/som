package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som/sombase"
	"net/url"
)

type URLExample struct {
	sombase.Node

	SomeURL      *url.URL
	SomeOtherURL url.URL
}
