package by

import (
	model "github.com/marcbinz/sdb/example/model"
	sort "github.com/marcbinz/sdb/lib/sort"
)

var User = newUser[model.User]("")

func newUser[T any](key string) user[T] {
	return user[T]{
		CreatedAt: sort.NewSort[T](keyed(key, "created_at")),
		Float32:   sort.NewSort[T](keyed(key, "float_32")),
		Float64:   sort.NewSort[T](keyed(key, "float_64")),
		ID:        sort.NewSort[T](keyed(key, "id")),
		Int:       sort.NewSort[T](keyed(key, "int")),
		Int32:     sort.NewSort[T](keyed(key, "int_32")),
		Int64:     sort.NewSort[T](keyed(key, "int_64")),
		String:    sort.NewString[T](keyed(key, "string")),
		UpdatedAt: sort.NewSort[T](keyed(key, "updated_at")),
		key:       key,
	}
}

type user[T any] struct {
	key       string
	ID        *sort.Sort[T]
	CreatedAt *sort.Sort[T]
	UpdatedAt *sort.Sort[T]
	String    *sort.String[T]
	Int       *sort.Sort[T]
	Int32     *sort.Sort[T]
	Int64     *sort.Sort[T]
	Float32   *sort.Sort[T]
	Float64   *sort.Sort[T]
}

func (n user[T]) MainGroup() group[T] {
	return newGroup[T](keyed(n.key, "main_group"))
}
