package model

import (
	"som.test/gen/som"
)

type Session struct {
	som.Node[som.UUID]
	som.Timestamps
	som.Expiry `som:"24h"`

	Token  string
	UserID string
}

// Ephemeral uses a very short TTL to exercise expiry filtering and purging.
type Ephemeral struct {
	som.Node[som.UUID]
	som.Expiry `som:"1s"`

	Label string
}
