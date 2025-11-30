//go:build embed

package lib

// SearchResult wraps a model with search metadata supporting multiple predicates.
// M is the model type.
type SearchResult[M any] struct {
	Model M
	// Scores maps predicate ref -> BM25 relevance score (0-1).
	Scores map[int]float64
	// Highlights maps predicate ref -> highlighted text with matched terms wrapped.
	Highlights map[int]string
}

// Score returns the score for the given predicate ref.
// If no ref is provided, returns the score for ref 0 (convenience for single-field search).
func (r SearchResult[M]) Score(ref ...int) float64 {
	if len(ref) == 0 {
		return r.Scores[0]
	}
	return r.Scores[ref[0]]
}

// Highlighted returns the highlighted text for the given predicate ref.
// If no ref is provided, returns the highlight for ref 0 (convenience for single-field search).
func (r SearchResult[M]) Highlighted(ref ...int) string {
	if len(ref) == 0 {
		return r.Highlights[0]
	}
	return r.Highlights[ref[0]]
}
