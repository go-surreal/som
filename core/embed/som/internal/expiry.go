//go:build embed

package internal

import "time"

// Expiry provides time-to-live functionality when embedded in a model.
// The duration is configured via the `som:"<duration>"` struct tag on the
// embed. The ExpiresAt field is set automatically by the database on creation.
type Expiry struct {
	expiresAt time.Time
}

// ExpiresAt returns the time at which this record expires.
func (e Expiry) ExpiresAt() time.Time {
	return e.expiresAt
}

func SetExpiresAt(e *Expiry, at time.Time) {
	e.expiresAt = at
}
