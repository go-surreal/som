//go:build embed

package lib

// Offset represents the start and end position of a matched term in the original text.
type Offset struct {
	Start int `cbor:"s"`
	End   int `cbor:"e"`
}

// SearchResult wraps a model with search metadata supporting multiple predicates.
// M is the model type.
type SearchResult[M any] struct {
	Model M
	// Scores contains all BM25 relevance scores (individual and combined).
	Scores []float64
	// Highlights maps predicate ref -> highlighted text with matched terms wrapped.
	Highlights map[int]string
	// Offsets maps predicate ref -> slice of position offsets for matched terms.
	Offsets map[int][]Offset
}

// Score returns the first score (convenience for single-field search).
func (r SearchResult[M]) Score() float64 {
	if len(r.Scores) == 0 {
		return 0
	}
	return r.Scores[0]
}

// Highlighted returns the highlighted text for the given predicate ref.
// If no ref is provided, returns the highlight for ref 0 (convenience for single-field search).
func (r SearchResult[M]) Highlighted(ref ...int) string {
	if len(ref) == 0 {
		return r.Highlights[0]
	}
	return r.Highlights[ref[0]]
}

// Offset returns the offsets for the given predicate ref.
// If no ref is provided, returns the offsets for ref 0 (convenience for single-field search).
func (r SearchResult[M]) Offset(ref ...int) []Offset {
	if len(ref) == 0 {
		return r.Offsets[0]
	}
	return r.Offsets[ref[0]]
}
