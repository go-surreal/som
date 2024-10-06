package field

import (
	"fmt"
)

type Bool struct {
	*BaseField
}

func (b *Bool) String() string {
	return fmt.Sprintf("%s: Bool", b.Name)
}
