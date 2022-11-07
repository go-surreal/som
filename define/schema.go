package define

type Schema struct {
	renders []Render
}

func (s *Schema) Render() []string {
	var statements []string
	for _, render := range s.renders {
		statements = append(statements, render.render())
	}
	return statements
}
