package field

import "fmt"

type String struct {
	*BaseField
}

func (s *String) String() string {
	return fmt.Sprintf("%s: String", s.Name)
}
