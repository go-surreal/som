package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som"
	"net/url"
)

type URLExample struct {
	som.Node

	Provider string `som:"unique(provider_account)"`
	Account  string `som:"unique(provider_account)"`

	SomeURL      *url.URL
	SomeOtherURL url.URL
}
