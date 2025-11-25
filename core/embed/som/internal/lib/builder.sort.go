//go:build embed

package lib

type SortBuilder struct {
	Field     string
	Order     SortOrder
	IsCollate bool
	IsNumeric bool
}

func (b *SortBuilder) render() string {
	out := b.Field + " "
	if b.IsCollate {
		out += "COLLATE "
	}
	if b.IsNumeric {
		out += "NUMERIC "
	}
	return out + string(b.Order)
}

// NewSortBuilder creates a new SortBuilder with the given field and order.
func NewSortBuilder(field string, order SortOrder) *SortBuilder {
	return &SortBuilder{
		Field: field,
		Order: order,
	}
}

type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)
