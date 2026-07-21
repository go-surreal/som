package model

import "som.test/gen/som"

type ChangefeedModel struct {
	som.Node[som.ULID] `som:"changefeed=1d"`
	som.Timestamps

	Name string
}
