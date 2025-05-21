package field

import "fmt"

type Field interface {
	fmt.Stringer

	Embedded() bool
}

type BaseField struct {
	Name string
}

func (f *BaseField) String() string {
	return fmt.Sprintf("%s: ?", f.Name)
}

func (f *BaseField) Embedded() bool {
	return f.Name == "" // TODO: would be true for types of pointer fields
}
