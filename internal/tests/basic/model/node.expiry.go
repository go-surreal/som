package model

import (
	"som.test/gen/som"
)

// Ephemeral uses a short expiry to exercise the ExpiresAt accessor, expiry
// filtering on read and background purging. Timestamps is embedded alongside
// Expiry to cover created_at and expires_at coexisting on the same table.
type Ephemeral struct {
	som.Node[som.UUID]
	som.Timestamps
	som.Expiry `som:"1s"`

	Label string
}
