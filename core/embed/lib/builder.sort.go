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

type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)
