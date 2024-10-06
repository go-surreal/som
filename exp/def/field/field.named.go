package field

import "fmt"

type Named struct {
	*BaseField

	Pkg      string
	TypeName string
}

func (n *Named) String() string {
	return fmt.Sprintf("%s: %s.%s", n.Name, n.Pkg, n.TypeName)
}
