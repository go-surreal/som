package other

// Amount is a type from another package that validates itself.
type Amount struct {
	Cents int
}

func (a Amount) Validate() error {
	return nil
}

// Plain has no validation at all.
type Plain struct {
	Value string
}
