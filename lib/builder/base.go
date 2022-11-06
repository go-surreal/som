package builder

type Block interface {
	render(*context) string
}
