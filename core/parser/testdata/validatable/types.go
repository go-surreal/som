package validatable

import (
	"errors"

	"github.com/go-surreal/som/core/parser/testdata/validatable/other"
)

// ValueRecv declares Validate with a value receiver.
type ValueRecv struct {
	Name string
}

func (v ValueRecv) Validate() error {
	return nil
}

// PtrRecv declares Validate with a pointer receiver.
type PtrRecv struct {
	Name string
}

func (v *PtrRecv) Validate() error {
	return nil
}

// NamedResult declares Validate with a named result.
type NamedResult struct{}

func (NamedResult) Validate() (err error) {
	return nil
}

// Enum is a named string type with its own validation.
type Enum string

func (e Enum) Validate() error {
	if e == "" {
		return errors.New("empty")
	}
	return nil
}

// Promoted inherits Validate from the embedded type.
type Promoted struct {
	ValueRecv
}

// Plain has no Validate method.
type Plain struct {
	Name string
}

// WrongResult returns something else than an error.
type WrongResult struct{}

func (WrongResult) Validate() bool {
	return true
}

// WrongParams takes an argument.
type WrongParams struct{}

func (WrongParams) Validate(strict bool) error {
	return nil
}

// TooManyResults returns more than one value.
type TooManyResults struct{}

func (TooManyResults) Validate() (bool, error) {
	return true, nil
}

// Fields holds one field per case, so that the detection can be run on types
// as they are encountered while parsing a model.
type Fields struct {
	Value    ValueRecv
	Ptr      *PtrRecv
	Enum     Enum
	Promoted Promoted
	Plain    Plain
	Imported other.Amount
	Foreign  other.Plain
	String   string
	Slice    []ValueRecv
}
