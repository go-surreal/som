package model

import (
	"errors"
	"strings"

	"som.test/gen/som"
)

// Validated exercises the generated write-time validation: the model itself,
// a named scalar, a nested struct, a slice of structs, optional values and a
// field that opted out.
type Validated struct {
	som.Node[som.ULID]

	Name string

	Contact    ContactAddress
	ContactPtr *ContactAddress
	Contacts   []ContactAddress

	Code    ProductCode
	CodePtr *ProductCode
	Codes   []ProductCode

	// Ignored holds a code that is not validated, even though its type
	// validates itself.
	Ignored ProductCode `som:"novalidate"`

	// Link points to another record, which is validated by its own repository
	// and must therefore not be checked here.
	Link *Validated

	Edges []ValidatedEdge
}

func (v *Validated) Validate() error {
	if v.Name == "" {
		return errors.New("name must not be empty")
	}

	return nil
}

// ContactAddress validates itself and holds a value that validates itself, so
// both checks have to run, the nested one first.
type ContactAddress struct {
	City  string
	Email ContactEmail
}

func (a ContactAddress) Validate() error {
	if a.City == "" {
		return errors.New("city must not be empty")
	}

	return nil
}

// ContactEmail is a named scalar validating itself.
type ContactEmail string

func (e ContactEmail) Validate() error {
	if !strings.Contains(string(e), "@") {
		return errors.New("not an email address")
	}

	return nil
}

// ProductCode is a named scalar with a pointer receiver.
type ProductCode string

func (c *ProductCode) Validate() error {
	if len(*c) != 4 {
		return errors.New("code must have four characters")
	}

	return nil
}

// ValidatedLog is a sink that validates its records before they are written.
type ValidatedLog struct {
	som.Sink

	Code ProductCode
}

// ValidatedEdge relates two validated models and validates itself.
type ValidatedEdge struct {
	som.Edge

	Validated Validated `som:"in"`
	AllTypes  AllTypes  `som:"out"`

	Note string
}

func (e ValidatedEdge) Validate() error {
	if e.Note == "" {
		return errors.New("note must not be empty")
	}

	return nil
}
