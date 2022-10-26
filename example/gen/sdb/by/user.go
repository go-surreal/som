package by

import (
	model "github.com/marcbinz/sdb/example/model"
	sort "github.com/marcbinz/sdb/lib/sort"
)

var User = newUser("")

func newUser(key string) user {
	return user{
		CreatedAt: sort.NewSort[model.User](key),
		Float32:   sort.NewSort[model.User](key),
		Float64:   sort.NewSort[model.User](key),
		ID:        sort.NewSort[model.User](key),
		Int:       sort.NewSort[model.User](key),
		Int32:     sort.NewSort[model.User](key),
		Int64:     sort.NewSort[model.User](key),
		String:    sort.NewString[model.User](key),
		UpdatedAt: sort.NewSort[model.User](key),
	}
}

type user struct {
	ID        *sort.Sort[model.User]
	CreatedAt *sort.Sort[model.User]
	UpdatedAt *sort.Sort[model.User]
	String    *sort.String[model.User]
	Int       *sort.Sort[model.User]
	Int32     *sort.Sort[model.User]
	Int64     *sort.Sort[model.User]
	Float32   *sort.Sort[model.User]
	Float64   *sort.Sort[model.User]
}
