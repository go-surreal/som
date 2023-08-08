package field

type Field interface {
	Embedded() bool
}

type BaseField struct {
	Name    string
	Pointer bool
}

func (f *BaseField) Embedded() bool {
	return f.Name == ""
}

type String struct {
	*BaseField
}
