package builder

type Sort struct {
	Field     string
	Order     SortOrder
	IsCollate bool
	IsNumeric bool
}

func (s *Sort) render() string {
	out := s.Field + " "
	if s.IsCollate {
		out += "COLLATE "
	}
	if s.IsNumeric {
		out += "NUMERIC "
	}
	return out + string(s.Order)
}

type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)
