//go:build embed

package lib

import (
	"strconv"
	"strings"
)

// SearchSort is implemented by types that can be used for sorting in search queries.
type SearchSort interface {
	SearchSort() *SortBuilder
}

type SortBuilder struct {
	Field     string
	Order     SortOrder
	IsCollate bool
	IsNumeric bool

	IsScore      bool
	ScoreRefs    []int
	ScoreMode    ScoreCombineMode
	ScoreWeights []float64
}

func (b *SortBuilder) SearchSort() *SortBuilder {
	return b
}

func (b *SortBuilder) render() string {
	if b.IsScore {
		refStrs := make([]string, len(b.ScoreRefs))
		for i, ref := range b.ScoreRefs {
			refStrs[i] = strconv.Itoa(ref)
		}
		return searchScorePrefix + strings.Join(refStrs, "_") + " " + string(b.Order)
	}

	// Due to a bug in SurrealDB when using ORDER BY with indexed fields,
	// we need to specifically SELECT all fields used for sorting with a
	// special alias to avoid issues for now.
	// see: https://github.com/surrealdb/surrealdb/issues/5588
	out := sortFieldPrefix + b.Field + " "

	if b.IsCollate {
		out += "COLLATE "
	}

	if b.IsNumeric {
		out += "NUMERIC "
	}

	return out + string(b.Order)
}

type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)
