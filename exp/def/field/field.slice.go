package field

import "fmt"

type Slice struct {
	*BaseField

	Elem Field
}

func (s *Slice) String() string {
	return fmt.Sprintf("%s: Slice(%s)", s.Name, s.Elem)
}
