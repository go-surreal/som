package field

import "fmt"

type Pointer struct {
	*BaseField

	Field Field
}

func (f *Pointer) String() string {
	return fmt.Sprintf("%s: Pointer(%s)", f.Name, f.Field)
}
