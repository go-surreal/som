package parser

type Node struct {
	Name       string
	Fields     []Field
	Timestamps bool
}

type Edge struct {
	Name       string
	In         Field
	Out        Field
	Fields     []Field
	Timestamps bool
}

type Struct struct {
	Name   string
	Fields []Field
}

type Enum struct {
	Name string
	Type string
}

type EnumValue struct {
	Enum     string
	Variable string
	Value    string
}
