package def

import "github.com/go-surreal/som/exp/def/field"

type Base struct {

	// Name is the name of the ?.
	Name string

	// Package is the fully qualified package name of the ?.
	// If it is empty, it is a built-in type.
	Package string

	// Exported defines whether the type is exported or not.
	Exported bool
}

type TypeParam struct {
	Name  string
	Field field.Field
}

func (tp *TypeParam) String() string {
	return tp.Name + ": " + tp.Field.String()
}
