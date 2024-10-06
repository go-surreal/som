package field

import (
	"fmt"
)

type Byte struct {
	*BaseField
}

func (b *Byte) String() string {
	return fmt.Sprintf("%s: Byte", b.Name)
}
