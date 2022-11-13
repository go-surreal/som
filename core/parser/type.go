package parser

type Node struct {
	Name   string
	Fields []Field
}

type Edge struct {
	Name   string
	In     Field
	Out    Field
	Fields []Field
}

type Struct struct {
	Name   string
	Fields []Field
}

type Enum struct {
	Name string
}

type EnumValue struct {
	Enum     string
	Variable string
	Value    string
}
