//go:build embed

package with

// Fetch_ has a suffix of "_" to prevent clashes with node names.
type Fetch_[T any] interface {
	fetch(T)
}

// FetchWithDeleted is implemented by fetch types for soft-delete models.
// It allows checking if the fetch should include soft-deleted records.
type FetchWithDeleted interface {
	IncludesDeleted() bool
	FetchField() string
}

func keyed[S ~string](base S, key string) string {
	if base == "" {
		return key
	}
	return string(base) + "." + key
}

func keyedStruct(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
