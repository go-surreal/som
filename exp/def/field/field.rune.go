package field

import (
	"fmt"
)

type Rune struct {
	*BaseField
}

func (r *Rune) String() string {
	return fmt.Sprintf("%s: Rune", r.Name)
}
