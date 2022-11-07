package statement

type Field struct {
	name string
}

func (f *Field) Render() string {
	return "DEFINE FIELD "
}
