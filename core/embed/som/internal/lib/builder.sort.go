//go:build embed

package lib

type SortBuilder struct {
	Field     string
	Order     SortOrder
	IsCollate bool
	IsNumeric bool
}

func (b *SortBuilder) render() string {
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
