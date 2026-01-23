package model

import "github.com/go-surreal/som/tests/basic/gen/som"

type ChangefeedModel struct {
	som.Node `som:"changefeed=1d"`
	som.Timestamps

	Name string
}
