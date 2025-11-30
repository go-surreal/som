//go:build embed

package lib

// ScoreSort represents score-based sorting for full-text search queries.
type ScoreSort struct {
	refs []int
	desc bool
}

// Score creates a new score-based sort by the given predicate refs.
// If multiple refs are provided, scores are summed.
func Score(refs ...int) *ScoreSort {
	return &ScoreSort{refs: refs, desc: true}
}

// Desc sets the sort order to descending (highest scores first).
// This is the default.
func (s *ScoreSort) Desc() *ScoreSort {
	s.desc = true
	return s
}

// Asc sets the sort order to ascending (lowest scores first).
func (s *ScoreSort) Asc() *ScoreSort {
	s.desc = false
	return s
}

// Refs returns the predicate refs used for scoring.
func (s *ScoreSort) Refs() []int {
	return s.refs
}

// IsDesc returns true if sorting descending.
func (s *ScoreSort) IsDesc() bool {
	return s.desc
}

func (s *ScoreSort) SearchSort() *SortBuilder {
	order := SortDesc
	if !s.desc {
		order = SortAsc
	}
	return &SortBuilder{
		IsScore:   true,
		ScoreRefs: s.refs,
		Order:     order,
	}
}
