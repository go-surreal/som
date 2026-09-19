package model

import (
	"som.test/gen/som"
)

// Constrained exercises the field-level constraints declared via som tags.
// The constraints are enforced by the database, so every field here maps to a
// DEFINE FIELD statement carrying an ASSERT clause.
type Constrained struct {
	som.Node[som.ULID]

	// Name must be between 3 and 64 characters long.
	Name string `som:"len=3..64"`

	// Code is a fixed-length identifier.
	Code string `som:"len=4"`

	// Nickname is optional, but constrained once present.
	Nickname *string `som:"len=2.."`

	// Age composes a tag constraint with the built-in range of the Go type.
	Age int `som:"min=0,max=130"`

	// Ratio only has a lower bound.
	Ratio float64 `som:"min=0"`

	// Tags may hold at most five entries, as len counts the items of a slice.
	Tags []string `som:"len=..5"`

	// Phone references a named constraint carrying a pattern.
	Phone string `som:"assert=phone_format"`

	// Title combines a tag constraint with a named raw expression.
	Title string `som:"len=1..80,assert=no_placeholder"`
}
