package builder

const (
	orderUnknown = iota
	orderWhere
	orderSort
)

type Block interface {
	render(*context) string
}
