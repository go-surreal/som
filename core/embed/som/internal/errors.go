//go:build embed

package internal

import (
	"errors"
)

// ErrInvalid is returned when a write is rejected because a value it carries is
// not allowed, no matter whether the rejection came from a Validate method or
// from a constraint the database enforces.
//
// It lives here rather than in the som package because the errors matching it
// do as well, and som imports this package, not the other way around. The som
// package re-exports it.
var ErrInvalid = errors.New("write rejected as invalid")
