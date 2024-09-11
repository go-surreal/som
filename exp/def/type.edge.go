package def

type Edge struct {
	*Struct
}

func (e *Edge) String() string {
	return e.describe("Edge")
}
