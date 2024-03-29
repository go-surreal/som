//go:build embed

package with

// Fetch_ has a suffix of "_" to prevent clashes with node names.
type Fetch_[T any] interface {
	fetch(T)
}

func keyed[S ~string](base S, key string) string {
	if base == "" {
		return key
	}
	return string(base) + "." + key
}
