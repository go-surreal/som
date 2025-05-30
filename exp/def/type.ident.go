package def

type Ident struct {
	*Base

	TypeParams []TypeParam
	Type       Type
}

func (s *Ident) String() string {
	return "Ident"
}

type Type struct {
}
