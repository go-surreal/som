package def

type Struct struct {
	Name       string
	TypeParams []TypeParam
	Fields     []Field
}

func (s *Struct) String() string {
	return s.describe("Struct")
}

func (s *Struct) describe(actualType string) string {
	out := s.Name + ": " + actualType + "(\n"

	if len(s.TypeParams) > 0 {
		out += "  TypeParams(\n"

		for _, tp := range s.TypeParams {
			out += "    " + tp.String() + "\n"
		}

		out += "  )\n"
	}

	if len(s.Fields) > 0 {
		out += "  Fields(\n"

		for _, f := range s.Fields {
			out += "    " + f.String() + "\n"
		}

		out += "  )\n"
	}

	out += ")"

	return out
}

type TypeParam struct {
	Name  string
	Field Field
}

func (tp *TypeParam) String() string {
	return tp.Name + ": " + tp.Field.String()
}
